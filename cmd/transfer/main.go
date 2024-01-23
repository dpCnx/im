package main

import (
	"flag"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"im/api/pb"
	"im/pkg/log"
)

var (
	flagConf string
)

func init() {
	flag.StringVar(&flagConf, "conf", "E:\\gowork\\im\\configs\\config.yaml", "config path, eg: -conf config.yaml")
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

	service, err := wireApp(bc.GetData(), log.NewLogger(bc.Log))
	if err != nil {
		panic(err)
	}
	service.MsgToMq()
	select {}
}
