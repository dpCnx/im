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

func NewData(msgRabbitmq *rabbitmq.MsgRabbitmq, logger log.Logger) *Data {
	return &Data{
		msgRabbitmq: msgRabbitmq,
		log:         log.NewHelper(logger),
	}
}

func (dt *Data) Consume() error {
	if err := dt.msgRabbitmq.Consume(func(bytes []byte) error {

		return nil
	}); err != nil {
		return err
	}
	return nil
}
