package repository

import (
	"errors"
	"sync"

	"github.com/rawmaterials223/MiniDouyin/util"
	"gorm.io/gorm"
)

// Videorelation DB
type VideoRelation struct {
	Id         int64 `gorm:"column:id"`
	FromUserId int64 `gorm:"column:from_user_id"`
	ToVideoId  int64 `gorm:"column:to_video_id"`
	IsLike     int   `gorm:"column:is_like"`
}

func (VideoRelation) TableName() string {
	return "videorelation"
}

type VideoRelationDao struct {
}

var videoRelationDao *VideoRelationDao
var videoRelationOnce sync.Once

func NewVideoRelationDaoInstance() *VideoRelationDao {
	videoRelationOnce.Do(
		func() {
			videoRelationDao = &VideoRelationDao{}
		})
	return videoRelationDao
}

// 查询：根据用户id和视频id查询点赞操作记录
func (*VideoRelationDao) QueryRelation(uid, vid int64) (int, error) {
	var relation VideoRelation

	// SQL: SELECT * FROM `videorelation` WHERE from_user_id = x
	// 		AND to_video_id = y ORDER BY id LIMIT 1;
	err := db.Where("from_user_id = ? AND to_video_id = ?", uid, vid).First(&relation).Error

	// 没有找到
	if errors.Is(err, gorm.ErrRecordNotFound) {
		util.Logger.Error("Query Video Relation ErrRecordNotFound")
		return 0, err
	}
	if err != nil {
		util.Logger.Error("Query videorelation error")
		return 0, err
	}

	return relation.IsLike, nil
}

// 更新：点赞记录变更，赞/取消赞
func (*VideoRelationDao) UpdateRelation(relation *VideoRelation) error {
	// SQL: UPDATE `videorelation` SET is_like = ?
	// 		WHERE from_user_id = x AND to_video_id = y;
	db.Model(&VideoRelation{}).
		Where("from_user_id = ? AND to_video_id = ?",
			relation.FromUserId,
			relation.ToVideoId).
		Update("is_like",
			relation.IsLike)

	util.Logger.Info("Update videorealtion success")
	return nil
}

// 插入：创建新的赞操作
func (*VideoRelationDao) CreateRelation(relation *VideoRelation) error {
	if err := db.Create(&relation).Error; err != nil {
		util.Logger.Error("Create Videorelation error")
		return err
	}

	return nil
}
