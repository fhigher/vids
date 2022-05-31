package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fhigher/vids/config"
	"github.com/fhigher/vids/core/jwtauth"
	"github.com/fhigher/vids/core/repo"
	"github.com/fhigher/vids/core/server"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/go-jsonrpc/auth"
	"github.com/gorilla/mux"
	"github.com/jinzhu/configor"
	"github.com/urfave/cli"
	"golang.org/x/xerrors"
)

var RunServerCmd = cli.Command{
	Name:  "run",
	Usage: "start server",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:   "repo",
			Usage:  "vids metadata home path",
			Value:  "~/.vids",
			EnvVar: repo.FsRepoEnv,
		},
		&cli.StringFlag{
			Name:     "config",
			Usage:    "specify config file path",
			Value:    "",
			Required: true,
			EnvVar:   config.ConfigPathEnv,
		},
		&cli.BoolFlag{
			Name:  "debug",
			Usage: "enable config debug mode",
		},
		&cli.StringFlag{
			Name:  "env",
			Usage: "which environment is going to run？ production,development or test",
			Value: config.Dev,
		},
		&cli.BoolFlag{
			Name:  "auto-reload",
			Usage: "auto reload config file",
		},
		&cli.IntFlag{
			Name:  "reload-intervals",
			Usage: "reload the configuration file every specified times, eg. 60s",
			Value: 60,
		},
	},
	Action: func(vctx *cli.Context) error {
		fs, err := repo.NewRepo(vctx.String("repo"))
		if err != nil {
			return err
		}
		err = fs.Init()
		if err != nil && err == repo.ErrRepoExists {
			log.Infof("repo exists: %s", fs.Path())
		} else {
			return err
		}

		if vctx.String("config") == "" {
			return xerrors.Errorf("server config path cannot empty.")
		}

		c := &configor.Config{
			Debug:              vctx.Bool("debug"),
			ENVPrefix:          "VIDS",
			Environment:        vctx.String("env"),
			AutoReload:         vctx.Bool("auto-reload"),
			AutoReloadInterval: time.Duration(vctx.Int("reload-intervals")) * time.Second,
			AutoReloadCallback: func(config interface{}) {
				log.Warning("server config reload")
			},
		}

		sc, err := config.InitConfig(vctx.String("config"), c)
		if nil != err {
			return err
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		rpcServer := jsonrpc.NewServer()

		s := server.NewVidsServer(&sc)
		rpcServer.Register("Vids", server.PermissionedAPI(s))

		mux := mux.NewRouter()
		mux.StrictSlash(true)
		mux.Handle("/rpc/v0", rpcServer)
		// 可以加一些metrics, promethues
		// mux.Handle("/metrics")

		j := jwtauth.NewJwtAuth("123456")
		ah := &auth.Handler{
			Verify: j.AuthVerify,
			Next:   mux.ServeHTTP,
			//Next: rpcServer.ServeHTTP,
		}

		srv := &http.Server{
			Handler: ah,
		}

		nl, err := net.Listen("tcp", sc.ServerListen+":"+sc.ServerPort)
		if err != nil {
			return err
		}

		go func() {
			log.Infof("setting up server endpoint at %s", nl.Addr().String())
			srv.Serve(nl)
		}()

		ch := make(chan os.Signal, 1)

		signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)

		sig := <-ch
		log.Warnw("received shutdown", "signal", sig)
		log.Warn("shutting down...")
		if err = srv.Shutdown(ctx); err != nil {
			return xerrors.Errorf("shutdown failed: %s", err)
		}

		log.Warn("Graceful shutdown successful")

		_ = log.Sync()

		return nil
	},
}
