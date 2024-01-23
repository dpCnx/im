//go:build wireinject
// +build wireinject

package main

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"im/api/pb"
	"im/internal/transfer/data"
	"im/internal/transfer/service"
	"im/pkg/rabbitmq"
)

func wireApp(*pb.Data, log.Logger) (*service.TransferService, error) {
	panic(wire.Build(service.ProviderSet, rabbitmq.ProviderSet, data.ProviderSet))
}
