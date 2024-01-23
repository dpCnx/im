package main

import (
	"flag"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	kLog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"im/api/pb"
	"im/internal/common"
	"im/pkg/log"
)

var (
	flagConf string
)

func init() {
	flag.StringVar(&flagConf, "conf", "E:\\gowork\\im\\configs\\config.yaml", "config path, eg: -conf config.yaml")
}

func newApp(gs *grpc.Server, hs *http.Server, logger kLog.Logger, r registry.Registrar) *kratos.App {
	return kratos.New(
		kratos.Name(common.Gateway),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			hs,
		),
		kratos.Registrar(r),
	)
}

func main() {
	flag.Parse()
	c := config.New(
		config.WithSource(
			file.NewSource(flagConf),
		),
	)
	defer c.Close()
	if err := c.Load(); err != nil {
		panic(err)
	}
	var bc pb.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	app, cleanup, err := wireApp(bc.Server.GetGateway(), bc.GetData(), log.NewLogger(bc.Log))
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
