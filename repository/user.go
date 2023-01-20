package repository

import (
	"sync"
	"time"
)

type User struct {
	Id         int64     `gorm:"column:id"`
	Name       string    `gorm:"column:name"`
	Token      string    `gorm:"column:token"`
	CreateTime time.Time `gorm:"column:create_time"`
}

type UserRelation struct {
	Id         int64 `gorm:"column:id"`
	FromUserId int64 `gorm:"column:from_user_id"`
	ToUserId   int64 `gorm:"column:to_user_id"`
	IsFollow   bool  `gorm:"column:is_follow"`
}

func (User) TableName() string {
	return "userinfo"
}

func (UserRelation) TableName() string {
	return "userrelation"
}

type UserDao struct {
}

var userDao *UserDao
var userOnce sync.Once

func NewUserDaoInstance() *UserDao {
	userOnce.Do(
		func() {
			userDao = &UserDao{}
		})
	return userDao
}

// check if the user exists
// user exist return nil
func (*UserDao) CheckUser(user *User) error {
	return nil
}

// insert the new user into table `userinfo`
func (*UserDao) CreateUser(user *User) error {
	return nil
}
