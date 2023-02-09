package repository

import (
	"sync"
	"time"
)

type Comment struct {
	Id         int64     `gorm:"column:id"`
	FromUserId int64     `gorm:"column:from_user_id"`
	ToVideoId  int64     `gorm:"column:to_video_id"`
	Content    string    `gorm:"column:content"`
	Status     int       `gorm:"column:status"`
	CreateTime time.Time `gorm:"column:create_time"`
	UpdateTime time.Time `gorm:"column:update_time"`
}

func (Comment) TableName() string {
	return "comment"
}

type CommentDao struct{}

var commentDao *CommentDao
var commentOnce sync.Once

func NewCommentDaoInstance() *CommentDao {
	commentOnce.Do(
		func() {
			commentDao = &CommentDao{}
		})

	return commentDao
}

func (*CommentDao) CalculateComment() (int, error) {
	var comment_count int

	return comment_count, nil
}
