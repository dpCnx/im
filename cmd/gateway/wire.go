//go:build wireinject
// +build wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"im/api/pb"
	"im/internal/gateway/conn"
	"im/internal/gateway/data"
	"im/internal/gateway/handler"
	"im/internal/gateway/server"
	"im/internal/gateway/service"
	"im/pkg/etcd"
)

func wireApp(*pb.GateWayConf, *pb.Data, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, service.ProviderSet,
		handler.ProviderSet, conn.ProviderSet, etcd.ProviderSet, data.ProviderSet,
		newApp))
}
