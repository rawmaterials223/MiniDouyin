package repository

import (
	"sync"
	"time"

	"github.com/rawmaterials223/MiniDouyin/util"
)

type Video struct {
	Id         int64     `gorm:"column:id"`
	UserId     int64     `gorm:"column:user_id"`
	PlayUrl    string    `gorm:"column:play_url"`
	CoverUrl   string    `gorm:"column:cover_url"`
	Title      string    `gorm:"column:title"`
	CreateTime time.Time `gorm:"column:create_time"`
}

func (Video) TableName() string {
	return "video"
}

type VideoDao struct{}

var videoDao *VideoDao
var videoOnce sync.Once

func NewVideoDaoInstance() *VideoDao {
	videoOnce.Do(
		func() {
			videoDao = &VideoDao{}
		})

	return videoDao
}

func (f *VideoDao) CreateVideo(video *Video) error {

	err := db.Create(&video).Error

	if err != nil {
		util.Logger.Error("Create Video Error: " + err.Error())
		return err
	}

	return nil
}
