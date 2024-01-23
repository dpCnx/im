package etcd

import (
	"time"

	"github.com/go-kratos/kratos/contrib/registry/etcd/v2"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/google/wire"
	etcdc "go.etcd.io/etcd/client/v3"
	ggrpc "google.golang.org/grpc"
	"im/api/pb"
)

var ProviderSet = wire.NewSet(NewRegister, NewDiscovery)

func NewRegister(config *pb.Data) registry.Registrar {

	client, err := etcdc.New(etcdc.Config{
		Endpoints:            []string{config.Etcd.Address},
		DialTimeout:          2 * time.Second,
		DialKeepAliveTime:    3 * time.Second, // 每3秒ping一次服务器
		DialKeepAliveTimeout: time.Second,     // 1秒没有返回则代表故障
		DialOptions:          []ggrpc.DialOption{ggrpc.WithBlock()},
	})
	if err != nil {
		panic(err)
	}

	return etcd.New(client)

}

func NewDiscovery(config *pb.Data) registry.Discovery {

	client, err := etcdc.New(etcdc.Config{
		Endpoints:            []string{config.Etcd.Address},
		DialTimeout:          2 * time.Second,
		DialKeepAliveTime:    3 * time.Second, // 每3秒ping一次服务器
		DialKeepAliveTimeout: time.Second,     // 1秒没有返回则代表故障
		DialOptions:          []ggrpc.DialOption{ggrpc.WithBlock()},
	})
	if err != nil {
		panic(err)
	}

	return etcd.New(client)

}
