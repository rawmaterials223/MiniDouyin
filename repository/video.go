package repository

import (
	"errors"
	"sync"
	"time"

	"github.com/rawmaterials223/MiniDouyin/util"
	"gorm.io/gorm"
)

// Video DB
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

// 查询：查询视频id是否存在，存在即返回视频信息
func (f *VideoDao) QueryVideoById(vid int64) (*Video, error) {
	var video Video

	// SQL: SELECT * FROM `video` where id = x
	err := db.Where("id = ?", vid).First(&video).Error

	// 没有找到记录
	if errors.Is(err, gorm.ErrRecordNotFound) {
		util.Logger.Error("Query Video ErrRecordNotFound")
		return nil, err
	}
	if err != nil {
		util.Logger.Error("Query Video Error")
		return nil, err
	}

	util.Logger.Info("Query Video success")
	return &video, nil
}

// 查询：查询用户uid的所有视频信息，包含点赞和评论数
func (f *VideoDao) QueryVideoByUid(uid int64) ([]VideoResult, error) {
	var videoResults []VideoResult

	/** SQL:
	SELECT t1.id, t1.play_url, t1.cover_url, t2.favorite_count, t1.title
	FROM `video` as t1
		left join
		(SELECT to_video_id as vid, count(from_user_id) as favorite_count
		 FROM `videorelation`
		 WHERE is_like = 1
		 GROUP BY to_video_id) as t2
	on t1.id = t2.vid
	where user_id = x;
	*/
	query := db.Table("videorelation").
		Select("to_video_id as vid, count(from_user_id) as favorite_count").
		Where("is_like = ?", 1).
		Group("to_video_id")
	result := db.Table("video").
		Select("video.id, video.play_url, video.cover_url, t2.favorite_count,0, video.title").
		Joins("left join (?) t2 on video.id = t2.vid", query).
		Where("user_id = ?", uid).
		Scan(&videoResults)

	if result.Error != nil {
		util.Logger.Error("Query Video Error")
		return nil, result.Error
	}
	// 查找多个对象
	// SQL: select * from video where user_id = x;
	//result = db.Where("user_id = ?", uid).Find(&videos)

	return videoResults, nil
}

func (f *VideoDao) CreateVideo(video *Video) error {

	err := db.Create(&video).Error

	if err != nil {
		util.Logger.Error("Create Video Error: " + err.Error())
		return err
	}

	return nil
}
