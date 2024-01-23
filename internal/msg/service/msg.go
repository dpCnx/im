package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"google.golang.org/protobuf/proto"
	"im/api/pb"
	"im/internal/msg/data"
)

var _ pb.MsgServer = (*MsgService)(nil)

type MsgService struct {
	pb.UnimplementedMsgServer
	log *log.Helper

	data *data.Data
}

func NewMsgService(logger log.Logger, data *data.Data) *MsgService {
	return &MsgService{log: log.NewHelper(logger), data: data}
}

func (m MsgService) SendMsg(ctx context.Context, req *pb.SendMsgReq) (*pb.SendMsgResp, error) {

	switch req.Msg.SessionType {
	case pb.SessionType_SINGLE_CHAT_TYPE:
		bytes, err := proto.Marshal(&pb.MessageToMq{
			PushUserId: req.Msg.RecvId,
			Message:    req.Msg,
		})
		if err != nil {
			return nil, pb.ErrorMsgSend(err.Error())
		}
		if err = m.data.PushMsg(bytes); err != nil {
			return nil, pb.ErrorMsgSend(err.Error())
		}

		if req.Msg.SendId != req.Msg.RecvId {
			bytes, err := proto.Marshal(&pb.MessageToMq{
				PushUserId: req.Msg.SendId,
				Message:    req.Msg,
			})
			if err != nil {
				return nil, pb.ErrorMsgSend(err.Error())
			}
			if err = m.data.PushMsg(bytes); err != nil {
				return nil, pb.ErrorMsgSend(err.Error())
			}
		}
	}
	return &pb.SendMsgResp{}, nil
}
