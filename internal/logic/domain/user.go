package domain

import "time"

type User struct {
	Id          uint64    `gorm:"column:id;primary_key;AUTO_INCREMENT;NOT NULL;comment:'自增主键'"`
	PhoneNumber string    `gorm:"column:phone_number;NOT NULL;comment:'手机号'"`
	Nickname    string    `gorm:"column:nickname;NOT NULL;comment:'昵称'"`
	Sex         int8      `gorm:"column:sex;NOT NULL;comment:'性别，0:未知；1:男；2:女'"`
	AvatarUrl   string    `gorm:"column:avatar_url;NOT NULL;comment:'用户头像链接'"`
	Extra       string    `gorm:"column:extra;NOT NULL;comment:'附加属性'"`
	CreateTime  time.Time `gorm:"column:create_time;default:CURRENT_TIMESTAMP;NOT NULL;comment:'创建时间'"`
	UpdateTime  time.Time `gorm:"column:update_time;default:CURRENT_TIMESTAMP;NOT NULL;comment:'更新时间'"`
}

func (u *User) TableName() string {
	return "user"
}
