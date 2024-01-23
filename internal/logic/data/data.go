package data

import (
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"im/api/pb"
	"im/internal/common"
)

var ProviderSet = wire.NewSet(
	NewData,
	NewMysql,
	NewRedis,
	NewLogicRepo,
	NewGetWayServiceClient,
)

type Data struct {
	db           *gorm.DB
	rds          *redis.Client
	getWayClient pb.GateWayClient
}

func NewData(db *gorm.DB, rds *redis.Client, getWayClient pb.GateWayClient) *Data {
	return &Data{
		db:           db,
		rds:          rds,
		getWayClient: getWayClient,
	}
}

func NewMysql(c *pb.Data) *gorm.DB {

	db, err := gorm.Open(mysql.Open(c.Database.Source), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}
	return db
}

func NewRedis(c *pb.Data) *redis.Client {

	rds := redis.NewClient(&redis.Options{
		Addr:     c.Redis.Address,
		Password: "",
	})

	_, err := rds.Ping(context.Background()).Result()

	if err != nil {
		panic(err)
	}
	return rds
}

func NewGetWayServiceClient(r registry.Discovery) pb.GateWayClient {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint(fmt.Sprintf("discovery:///%s", common.Gateway)),
		grpc.WithDiscovery(r),
		grpc.WithMiddleware(
			recovery.Recovery(),
		),
	)
	if err != nil {
		panic(err)
	}
	return pb.NewGateWayClient(conn)
}
