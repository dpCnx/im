package data

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"im/pkg/rabbitmq"
)

var ProviderSet = wire.NewSet(
	NewData,
)

type Data struct {
	msgRabbitmq *rabbitmq.MsgRabbitmq
	log         *log.Helper
}

func NewData(msgRabbitmq *rabbitmq.MsgRabbitmq, logger log.Logger) (*Data, error) {
	tmp := &Data{
		msgRabbitmq: msgRabbitmq,
		log:         log.NewHelper(logger),
	}
	if err := tmp.CreateAndBind(); err != nil {
		return nil, err
	}
	tmp.Confirm()
	tmp.Return()

	return tmp, nil
}

func (dt *Data) CreateAndBind() error {
	if err := dt.msgRabbitmq.CreateExchange(); err != nil {
		return err
	}
	_, err := dt.msgRabbitmq.CreateQueue()
	if err != nil {
		return err
	}
	if err = dt.msgRabbitmq.BindQueue(); err != nil {
		return err
	}
	return nil
}

func (dt *Data) Confirm() {
	_ = dt.msgRabbitmq.Confirm(func(deliveryTag uint64) {
		dt.log.Error(deliveryTag)
	})
}

func (dt *Data) Return() {
	_ = dt.msgRabbitmq.Return(func(res []byte) {
		dt.log.Error(string(res))
	})
}

func (dt *Data) PushMsg(msg []byte) error {
	return dt.msgRabbitmq.PushMsg(msg)
}
