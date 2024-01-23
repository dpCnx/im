//go:build wireinject
// +build wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"im/api/pb"
	"im/internal/msg/data"
	"im/internal/msg/server"
	"im/internal/msg/service"
	"im/pkg/etcd"
	"im/pkg/rabbitmq"
)

func wireApp(*pb.MsgConfig, *pb.Data, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, service.ProviderSet, etcd.ProviderSet, rabbitmq.ProviderSet, data.ProviderSet, newApp))
}
