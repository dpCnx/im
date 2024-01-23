//go:build wireinject
// +build wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"im/api/pb"
	"im/internal/logic/data"
	"im/internal/logic/server"
	"im/internal/logic/service"
	"im/pkg/etcd"
)

func wireApp(*pb.LogicConfig, *pb.Data, log.Logger) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, data.ProviderSet, service.ProviderSet, etcd.ProviderSet, newApp))
}
