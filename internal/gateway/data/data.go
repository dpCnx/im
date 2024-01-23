package data

import (
	"context"
	"fmt"

	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/google/wire"
	"im/api/pb"
	"im/internal/common"
)

var ProviderSet = wire.NewSet(
	NewData,
	NewMsgServiceClient,
)

type Data struct {
	MsgClient pb.MsgClient
}

func NewData(msgClient pb.MsgClient) *Data {
	return &Data{MsgClient: msgClient}
}

func NewMsgServiceClient(r registry.Discovery) pb.MsgClient {
	conn, err := grpc.DialInsecure(
		context.Background(),
		grpc.WithEndpoint(fmt.Sprintf("discovery:///%s", common.Msg)),
		grpc.WithDiscovery(r),
		grpc.WithMiddleware(
			recovery.Recovery(),
		),
	)
	if err != nil {
		panic(err)
	}
	return pb.NewMsgClient(conn)
}
