package repository

import (
	"errors"
	"strconv"
	"sync"

	"github.com/rawmaterials223/MiniDouyin/util"
	"gorm.io/gorm"
)

func (User) TableName() string {
	return "userinfo"
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
		util.Logger.Error("Query User ErrRecordNotFound")
		return nil, err
	}
	if err != nil {
		util.Logger.Error("Query User Error: " + err.Error())
		return nil, err
	}

	util.Logger.Info("Query User userId = " + strconv.FormatInt(user.Id, 10) + " , usertoken = " + user.Token)

	return &user, nil
}

// check if the user exists
// user exist return (&user, nil)
// doesn't exist return (nil, nil)
// other error return(nil, err)
func (*UserDao) QueryUserByIdToken(uid int64, token string) (*User, error) {
	var user User

	// SQL: SELECT * FROM `userinfo` WHERE id = x and token = y ORDER BY id LIMIT 1;
	err := db.Where("id = ? AND token = ?", uid, token).First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		util.Logger.Error("Query User ErrRecordNotFound")
		return nil, err
	}
	if err != nil {
		util.Logger.Error("Query User Error: " + err.Error())
		return nil, err
	}

	return &user, nil
}

func (*UserDao) QueryUserByToken(token string) (*User, error) {
	var user User

	err := db.Where("token = ?", token).First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		util.Logger.Error("Query User ErrRecordNotFound")
		return nil, err
	}
	if err != nil {
		util.Logger.Error("Query User Error: " + err.Error())
		return nil, err
	}

	return &user, nil
}

// insert the new user into table `userinfo`
func (*UserDao) CreateUser(user *User) error {

	err := db.Create(&user).Error

	if err != nil {
		util.Logger.Error("Create User Error: " + err.Error())
		return err
	}

	return nil
}
