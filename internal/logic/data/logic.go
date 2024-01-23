package data

import (
	"context"

	"github.com/go-kratos/kratos/v2/errors"
	"gorm.io/gorm"
	"im/api/pb"
	"im/internal/common"
)

type LogicRepo struct {
	data *Data
}

func NewLogicRepo(data *Data) *LogicRepo {
	return &LogicRepo{
		data: data,
	}
}

func (lr *LogicRepo) SendMessage(ctx context.Context, userId int64, msg *pb.Message) (*pb.PushMsgResp, error) {

	incr, err := lr.SeqIncr(common.SeqObjectTypeUser, userId)
	if err != nil {
		return nil, err
	}
	msg.Seq = incr
	return lr.data.getWayClient.PushMsg(ctx, &pb.PushMsgReq{
		PushToUserIdList: []int64{userId},
		Message:          msg,
	})
}

func (lr *LogicRepo) SeqIncr(objectType int, objectId int64) (int64, error) {
	tx := lr.data.db.Begin()

	var seq int64
	err := tx.Raw("select seq from seq where object_type = ? and object_id = ? for update", objectType, objectId).Row().Scan(&seq)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return 0, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if err = tx.Exec("insert into seq (object_type,object_id,seq) values (?,?,?)", objectType, objectId, seq+1).Error; err != nil {
			tx.Rollback()
			return 0, err
		}
	} else {
		if err = tx.Exec("update seq set seq = seq + 1 where object_type = ? and object_id = ?", objectType, objectId).Error; err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	tx.Commit()
	return seq + 1, nil
}
