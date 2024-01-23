package service

import (
	"github.com/go-kratos/kratos/v2/log"
	"im/internal/transfer/data"
)

type TransferService struct {
	log  *log.Helper
	data *data.Data
}

func NewPushService(logger log.Logger, data *data.Data) *TransferService {
	return &TransferService{
		log:  log.NewHelper(logger),
		data: data,
	}
}

func (ts *TransferService) MsgToMq() {
	ts.data.Consume()
}
