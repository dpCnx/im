package handler

import (
	"context"

	"google.golang.org/protobuf/proto"
	"im/api/pb"
	"im/internal/common"
	"im/internal/gateway/conn"
)

func (ws *WsHandler) msgParse(c *conn.Conn, bytes []byte) {
	var input pb.Input
	if err := proto.Unmarshal(bytes, &input); err != nil {
		ws.log.Error("proto unmarshal", err)
		if err = c.Send(pb.PackageType_PT_MESSAGE, common.ErrServer, err.Error(), nil); err != nil {
			ws.log.Error("send err", err)
			return
		}
		return
	}
	ws.log.Info("ws msg", input)

	var err error
	switch input.Type {
	case pb.PackageType_PT_SYNC:
		err = ws.Sync(&input)
	case pb.PackageType_PT_HEARTBEAT:
		err = ws.Heartbeat(&input)
	case pb.PackageType_PT_MESSAGE:
		err = ws.Message(&input)
	default:
		ws.log.Info("handler switch other")
	}
	if err != nil {
		if err = c.Send(pb.PackageType_PT_MESSAGE, common.ErrServer, err.Error(), nil); err != nil {
			ws.log.Error("send err", err)
			return
		}
	}
}

func (ws *WsHandler) Sync(input *pb.Input) error {
	return nil
}

func (ws *WsHandler) Heartbeat(input *pb.Input) error {
	return nil
}

func (ws *WsHandler) Message(input *pb.Input) error {
	var message pb.Message
	if err := proto.Unmarshal(input.Data, &message); err != nil {
		return err
	}
	_, err := ws.data.MsgClient.SendMsg(context.Background(), &pb.SendMsgReq{Msg: &message})
	if err != nil {
		return err
	}
	return nil
}
