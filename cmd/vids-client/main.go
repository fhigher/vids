package main

import (
	"os"

	logging "github.com/ipfs/go-log/v2"
	"github.com/urfave/cli"
)

var log = logging.Logger("vids-client")

func main() {

	SetupLogLevels()

	app := cli.NewApp()
	app.Name = "vids-client"
	app.Usage = "the client of vids"
	app.Version = "v0.0.1"
	app.Authors = []cli.Author{
		{
			Name:  "Jerry",
			Email: "1113821597@qq.com",
		},
	}
	app.EnableBashCompletion = true
	app.Commands = []cli.Command{
		ImportMysqlCmd,
	}

	if err := app.Run(os.Args); nil != err {
		log.Error(err)
		os.Exit(1)
	}
}

func SetupLogLevels() {
	if _, set := os.LookupEnv("GOLOG_LOG_LEVEL"); !set {
		_ = logging.SetLogLevel("*", "DEBUG")
	}
}
