package service

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/rawmaterials223/MiniDouyin/repository"
	"github.com/rawmaterials223/MiniDouyin/util"
)

func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

type ResponseError struct {
	Status  int
	Message string
}

func (r *ResponseError) Error() string {
	return fmt.Sprintf("%v", r.Message)
}

func Register(username, password string) (int64, string, error) {
	return NewRegisterFlow(username, password).DoRegister()
}

func NewRegisterFlow(username, password string) *UserFlow {
	return &UserFlow{
		username: username,
		password: password,
	}
}

func Login(username, password string) (int64, string, error) {
	return NewLoginFlow(username, password).DoLogin()
}

func NewLoginFlow(username, password string) *UserFlow {
	return &UserFlow{
		username: username,
		password: password,
	}
}

func Info(uid, token string) (UserResultType, error) {
	return NewUserInfoFlow(uid, token).Do()
}

func NewUserInfoFlow(uid, token string) *UserInfoFlow {
	s, _ := strconv.ParseInt(uid, 10, 64)
	return &UserInfoFlow{
		userId: s,
		token:  token,
	}
}

type UserResultType *repository.UserResult

type UserFlow struct {
	username string
	password string

	userId int64
	token  string
}

type UserInfoFlow struct {
	userId int64
	token  string

	UserResult UserResultType
}

func (f *UserFlow) DoRegister() (int64, string, error) {

	// check if user existed
	exist, _ := f.CheckUserByNameToken()
	if exist {
		return 0, "", &ResponseError{1, "register error: user existed"}
	}

	// register
	if err := f.CreateUser(); err != nil {
		return 0, "", err
	}
	// 注册成功后返回ID和token
	return f.userId, f.token, nil
}

func (f *UserFlow) DoLogin() (int64, string, error) {

	// check if user existed
	exist, _ := f.CheckUserByNameToken()
	if !exist {
		return 0, "", &ResponseError{1, "register error: user doesn't exist"}
	}

	// 登录成功后返回ID和token
	return f.userId, f.token, nil
}

// 功能：检查用户是否存在，通过昵称和密码
// 用户存在-(true, nil)
// 用户不存在-(false, err)
func (f *UserFlow) CheckUserByNameToken() (bool, error) {

	// receive result from repository layer
	// user exist return (&user, nil)
	// doesn't exist return (nil, nil)
	// other error return(nil, err)
	user, err := repository.NewUserDaoInstance().QueryUserByNameToken(f.username, f.username+f.password)
	// 用户不存在或查找出错
	if err != nil {
		return false, err
	}

	util.Logger.Info("QueryUserByNameToken success")

	f.userId = user.Id
	f.token = user.Token

	return true, nil
}

func (f *UserFlow) CreateUser() error {

	user := &repository.User{
		Name:       f.username,
		Token:      f.username + f.password, // TODO: token 计算方式待定
		CreateTime: time.Now(),
	}

	// receive result from repository layer
	// if insert failed, return error message
	if err := repository.NewUserDaoInstance().CreateUser(user); err != nil {
		return err
	}

	f.userId = user.Id
	f.token = user.Token

	fmt.Printf("CreateUser userId = %d", f.userId)

	return nil
}

func (f *UserInfoFlow) Do() (*repository.UserResult, error) {

	// 检查用户是否存在
	exist, _ := f.CheckUserByIdToken()
	if !exist {
		return nil, &ResponseError{1, "user doesn't exist"}
	}

	// 计算关注数follow_count，粉丝数follower_count，是否关注is_follow
	if err := f.QueryUserRelationById(); err != nil {
		return nil, err
	}

	return f.UserResult, nil
}

func (f *UserInfoFlow) CheckUserByIdToken() (bool, error) {

	_, err := repository.NewUserDaoInstance().QueryUserByIdToken(f.userId, f.token)

	if err != nil {
		return false, err
	}

	util.Logger.Info("QueryUserByIdToken success")

	return true, nil
}

func (f *UserInfoFlow) QueryUserRelationById() error {
	userResult, err := repository.NewRelationDaoInstance().Calcualte(f.userId)

	if err != nil {
		return err
	}

	isFollow, _ := repository.NewRelationDaoInstance().QueryRelation(f.userId, f.userId)

	f.UserResult = &userResult
	if isFollow == DoActionType {
		f.UserResult.IsFollow = true
	} else {
		f.UserResult.IsFollow = false
	}
	util.Logger.Info("QueryUserResult success")

	return nil
}
