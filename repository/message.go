package repository

import (
	"sync"
	"time"
)

type Message struct {
	Id         int64     `gorm:"column:id"`
	FromUserId int64     `gorm:"column:from_user_id"`
	ToUserId   int64     `gorm:"column:to_user_id"`
	Content    string    `gorm:"column:content"`
	CreateTime time.Time `gorm:"column:create_time"`
}

func (Message) TableName() string {
	return "message"
}

type MessageDao struct {
}

var messageDao *MessageDao
var messageOnce sync.Once

func NewMessageDaoInstance() *MessageDao {
	messageOnce.Do(
		func() {
			messageDao = &MessageDao{}
		})
	return messageDao
}
