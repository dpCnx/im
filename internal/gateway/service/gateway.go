package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"im/api/pb"
	"im/internal/common"
	"im/internal/gateway/conn"
)

var _ pb.GateWayServer = (*GateWayService)(nil)

type GateWayService struct {
	pb.UnimplementedGateWayServer

	log         *log.Helper
	connManager *conn.Manager
}

func NewGateWayService(connManager *conn.Manager, logger log.Logger) *GateWayService {
	return &GateWayService{log: log.NewHelper(logger), connManager: connManager}
}

func (g GateWayService) PushMsg(ctx context.Context, req *pb.PushMsgReq) (*pb.PushMsgResp, error) {
	for _, pushUserId := range req.PushToUserIdList {
		cons := g.connManager.GetUserAllCons(pushUserId)
		for _, conns := range cons {
			for _, c := range conns {
				if err := c.Send(pb.PackageType_PT_MESSAGE, common.Success, "", req.Message); err != nil {
					return nil, pb.ErrorMsgSend(err.Error())
				}
			}
		}
	}
	return &pb.PushMsgResp{}, nil
}
