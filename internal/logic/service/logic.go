package service

import (
	"context"

	"im/api/pb"
	"im/internal/logic/data"
)

var _ pb.LogicServer = (*LogicService)(nil)

type LogicService struct {
	pb.UnimplementedLogicServer

	logicRepo *data.LogicRepo
}

func NewLogicService(logicRepo *data.LogicRepo) *LogicService {
	return &LogicService{logicRepo: logicRepo}
}

func (l LogicService) AddFriend(ctx context.Context, req *pb.AddFriendReq) (*pb.AddFriendResp, error) {
	// TODO implement me
	panic("implement me")
}
