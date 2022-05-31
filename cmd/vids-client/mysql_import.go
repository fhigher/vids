package main

import (
	"context"

	mc "github.com/fhigher/vids/datasource/mysql"
	"github.com/urfave/cli"
	"golang.org/x/xerrors"
)

var ImportMysqlCmd = cli.Command{
	Name:  "msql-import",
	Usage: "import mysql data source",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "host",
			Usage:    "mysql server host addr",
			Required: true,
		},
		&cli.StringFlag{
			Name:  "port",
			Usage: "mysql server port",
			Value: "3306",
		},
		&cli.StringFlag{
			Name:  "user",
			Usage: "mysql server user name",
			Value: "root",
		},
		&cli.StringFlag{
			Name:  "pass",
			Usage: "mysql server password of user",
		},
		&cli.StringFlag{
			Name:     "dbname",
			Usage:    "mysql server datebase name",
			Required: true,
		},
		&cli.StringFlag{
			Name:  "tbnames",
			Usage: "which tables will be import. eg:'buy_log,user_order",
		},
		&cli.IntFlag{
			Name:  "max-conn",
			Usage: "set MaxOpenConns",
			Value: 100,
		},
	},
	Action: func(vctx *cli.Context) error {
		if !vctx.IsSet("host") {
			return xerrors.Errorf("mysql server host must be set")
		}
		if !vctx.IsSet("dbname") {
			return xerrors.Errorf("mysql server dbname must be set")
		}
		if !vctx.IsSet("tbnames") {
			log.Warn("no tbnames specify. use all tables")
		}

		dbInfo := mc.NewDbInfo(vctx.String("host"), vctx.String("port"), vctx.String("user"),
			vctx.String("pass"), vctx.String("dbname"), vctx.Int("max-conn"))

		dbInfo.InitDB()

		deal := mc.NewDealMysqlData(dbInfo, vctx.String("tbnames"))
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		err := deal.Start(ctx)
		if err != nil {
			log.Error(err)
		}

		return nil
	},
}
