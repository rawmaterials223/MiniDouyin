package repository

import (
	"sync"

	"github.com/rawmaterials223/MiniDouyin/util"
)

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
