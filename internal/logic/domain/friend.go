package domain

import "time"

type Friend struct {
	Id         uint64    `gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL;comment:'自增主键'"`
	UserId     uint64    `gorm:"column:user_id;NOT NULL;comment:'用户id'"`
	FriendId   uint64    `gorm:"column:friend_id;NOT NULL;comment:'好友id'"`
	Remarks    string    `gorm:"column:remarks;NOT NULL;comment:'备注'"`
	Extra      string    `gorm:"column:extra;NOT NULL;comment:'附加属性'"`
	CreateTime time.Time `gorm:"column:create_time;default:CURRENT_TIMESTAMP;NOT NULL;comment:'创建时间'"`
	UpdateTime time.Time `gorm:"column:update_time;default:CURRENT_TIMESTAMP;NOT NULL;comment:'更新时间'"`
}

func (f *Friend) TableName() string {
	return "friend"
}
