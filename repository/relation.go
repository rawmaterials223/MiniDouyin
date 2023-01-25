package repository

import (
	"errors"
	"fmt"
	"sync"

	"gorm.io/gorm"
)

type UserRelation struct {
	Id         int64 `gorm:"column:id"`
	FromUserId int64 `gorm:"column:from_user_id"`
	ToUserId   int64 `gorm:"column:to_user_id"`
	IsFollow   int   `gorm:"column:is_follow"`
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
		return 0, err
	}
	if err != nil {
		return 0, err
	}

	fmt.Printf("Query User userId = %d, usertoken = %v ", user.Id, user.Token)

	return user.Id, nil
}

func (*RelationDao) CreateAction(relation *UserRelation) error {

	return nil
}
