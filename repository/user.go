package repository

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"gorm.io/gorm"
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
// user exist return (&user, nil)
// doesn't exist return (nil, nil)
// other error return(nil, err)
func (*UserDao) QueryUserByNameToken(username string, token string) (*User, error) {
	var user User

	// SQL: SELECT * FROM `userinfo` WHERE name = x and token = y ORDER BY id LIMIT 1;
	err := db.Where("name = ? AND token = ?", username, token).First(&user).Error

	// 没有找到记录,ErrRecordNotFound与First配合使用
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	fmt.Printf("Query User userId = %d, usertoken = %v ", user.Id, user.Token)

	return &user, nil
}

// insert the new user into table `userinfo`
func (*UserDao) CreateUser(user *User) error {

	err := db.Create(&user).Error

	if err != nil {
		return err
	}

	return nil
}
