package domain

import "time"

type Seq struct {
	Id         uint64    `gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL;comment:'自增主键'"`
	ObjectType int8      `gorm:"column:object_type;NOT NULL;comment:'对象类型,1:用户；2：群组'"`
	ObjectId   uint64    `gorm:"column:object_id;NOT NULL;comment:'对象id'"`
	Seq        uint64    `gorm:"column:seq;NOT NULL;comment:'序列号'"`
	CreateTime time.Time `gorm:"column:create_time;default:CURRENT_TIMESTAMP;NOT NULL;comment:'创建时间'"`
	UpdateTime time.Time `gorm:"column:update_time;default:CURRENT_TIMESTAMP;NOT NULL;comment:'更新时间'"`
}

func (s *Seq) TableName() string {
	return "seq"
}
