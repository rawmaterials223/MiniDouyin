package repository

import (
	"errors"
	"strconv"
	"sync"

	"github.com/rawmaterials223/MiniDouyin/util"
	"gorm.io/gorm"
)

type UserRelation struct {
	Id         int64 `gorm:"column:id"`
	FromUserId int64 `gorm:"column:from_user_id"`
	ToUserId   int64 `gorm:"column:to_user_id"`
	IsFollow   int   `gorm:"column:is_follow"`
}

type UserRelationCount struct {
	FollowCount   int  `json:"follow_count"`
	FollowerCount int  `json:"follower_count"`
	IsFollow      bool `json:"is_follow"`
}

func (UserRelation) TableName() string {
	return "userrelation"
}

type RelationDao struct {
}

var relationDao *RelationDao
var relationOnce sync.Once

func NewRelationDaoInstance() *RelationDao {
	relationOnce.Do(
		func() {
			relationDao = &RelationDao{}
		})
	return relationDao
}

func (*RelationDao) QueryUserByToken(token string) (int64, error) {
	var user User

	// SQL: SELECT * FROM `userinfo` WHERE token = y ORDER BY id LIMIT 1;
	err := db.Where("token = ?", token).First(&user).Error

	// 没有找到记录,ErrRecordNotFound与First配合使用
	if errors.Is(err, gorm.ErrRecordNotFound) {
		util.Logger.Error("Query User ErrRecordNotFound")
		return 0, err
	}
	if err != nil {
		util.Logger.Error("Query User Error: " + err.Error())
		return 0, err
	}

	util.Logger.Info("Query User userId = " + strconv.FormatInt(user.Id, 10) + " , usertoken = " + user.Token)

	return user.Id, nil
}

// 查询：用户关系记录
func (*RelationDao) QueryRelation(from_id, to_id int64) (int, error) {
	var relation UserRelation

	// SQL: SELECT * FROM `userrelation` WHERE from_user_id = x AND to_user_id = y ORDER BY id LIMIT 1;
	err := db.Where("from_user_id = ? AND to_user_id = ?", from_id, to_id).First(&relation).Error

	// 没有找到
	if errors.Is(err, gorm.ErrRecordNotFound) {
		util.Logger.Error("Query Relation ErrRecordNotFound")
		return 0, err
	}
	if err != nil {
		util.Logger.Error("Query userrelation Error: " + err.Error())
		return 0, err
	}

	return relation.IsFollow, nil
}

// 更新：用户关系变更，关注/取消关注
func (*RelationDao) UpdateRelation(relation *UserRelation) error {

	// 不能用save，缺少where条件，会添加转为添加新数据
	// db.Save(&relation) -> UPDATE userrelation SET from_user_id = x, to_user_id = y, is_follow = z;

	// SQL: UPDATE userrelation SET is_follow = x WHERE from_user_id = y AND to_user_id = z;
	db.Model(&UserRelation{}).
		Where("from_user_id = ? AND to_user_id = ?",
			relation.FromUserId,
			relation.ToUserId).
		Update("is_follow",
			relation.IsFollow)

	util.Logger.Info("Update UserRelation Success")
	return nil
}

// 插入：用户关系表插入新数据
func (*RelationDao) CreateRelation(relation *UserRelation) error {

	if err := db.Create(&relation).Error; err != nil {
		util.Logger.Error("Create UserRelation Error: " + err.Error())
		return err
	}

	return nil
}

// 查询：计算用户的关注数和粉丝数
func (*RelationDao) CalculateRelation(uid int64) (int64, int64, error) {
	var follow_count int64
	var follower_count int64

	// follow_count
	// SQL: select count(*) from `userrelation` where from_user_id = x and is_follow = 1;
	db.Model(&UserRelation{}).Where("from_user_id = ? AND is_follow = ?", uid, 1).Count(&follow_count)

	// follower_count
	// SQL: select count(*) from `userrelation` where to_user_id = x and is follow = 1;
	db.Model(&UserRelation{}).Where("to_user_id = ? AND is_follow = ?", uid, 1).Count(&follower_count)

	return follow_count, follower_count, nil
}
