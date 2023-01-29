package repository

import (
	"sync"
	"time"
)

type Comment struct {
	Id         int64     `gorm:"column:id"`
	VideoId    int64     `gorm:"column:video_id"`
	UserId     int64     `gorm:"column:user_id"`
	Content    string    `gorm:"column:content"`
	CreateTime time.Time `gorm:"column:create_time"`
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
