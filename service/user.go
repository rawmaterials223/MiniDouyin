package service

import (
	"fmt"
	"time"

	"github.com/rawmaterials223/MiniDouyin/repository"
)

type RegisterError struct {
	Status  int
	Message string
}

func (r *RegisterError) Error() string {
	return fmt.Sprintf("%v", r.Message)
}

func Register(username string, password string) (int64, string, error) {
	return NewRegisterFlow(username, password).DoRegister()
}

func NewRegisterFlow(username string, password string) *UserFlow {
	return &UserFlow{
		username: username,
		password: password,
	}
}

func Login(username string, password string) (string, error) {
	return NewLoginFlow(username, password).DoLogin()
}

func NewLoginFlow(username string, password string) *UserFlow {
	return &UserFlow{
		username: username,
		password: password,
	}
}

type UserFlow struct {
	username string
	password string

	userId int64
	token  string
}

func (f *UserFlow) DoRegister() (int64, string, error) {

	// check if user existed
	if err := f.IsExistedUser(); err == nil {
		return 0, "", &RegisterError{1, "user existed"}
	}

	// register
	if err := f.CreateUser(); err != nil {
		return 0, "", err
	}

	return f.userId, f.token, &RegisterError{Status: 0}
}

// TODO: check if user existed
// user exist return nil
// else return RegisterError
func (f *UserFlow) IsExistedUser() error {

	user := &repository.User{
		Name:  f.username,
		Token: f.token,
	}

	// receive result from repository layer
	if err := repository.NewUserDaoInstance().CheckUser(user); err != nil {
		return &RegisterError{Status: 1, Message: "User doesn't exist"}
	}

	// TODO: token的获取方式待定
	f.token = user.Token

	return nil
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
		return &RegisterError{Status: 1, Message: "failed"} // TODO: err message待定
	}

	f.userId = user.Id
	f.token = user.Token

	return nil
}

func (f *UserFlow) DoLogin() (string, error) {
	if err := f.IsExistedUser(); err != nil {
		return "", &RegisterError{1, "user doesn't exist"}
	}

	return f.token, nil
}
