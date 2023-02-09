package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rawmaterials223/MiniDouyin/service"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin
var usersLoginInfo = map[string]User{
	"zhangleidouyin": {
		Id:            1,
		Name:          "zhanglei",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
	},
}

//var userIdSequence = int64(1)

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User service.UserResultDemo `json:"user"`
}

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	// receive the result from service layer
	// return userId AND token
	// if register failed, StatusMsg can be taken from err.Error()
	// 注册成功返回id, token, err/nil
	// 注册失败返回0, "", err/nil
	userId, token, err := service.Register(username, password)

	if userId == 0 && token == "" && err != nil {
		c.JSON(
			http.StatusOK,
			UserLoginResponse{
				Response: Response{
					StatusCode: 1,
					StatusMsg:  err.Error(),
				},
			})
	} else {
		c.JSON(
			http.StatusOK,
			UserLoginResponse{
				Response: Response{
					StatusCode: 0,
				},
				UserId: userId,
				Token:  token,
			})
	}
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	//token := username + password

	// receive the result from service layer
	// return userid AND token
	// if user doesn't exist return err, else return nil
	userid, token, err := service.Login(username, password)

	if userid == 0 && token == "" && err != nil {
		c.JSON(
			http.StatusOK,
			UserLoginResponse{
				Response: Response{
					StatusCode: 1,
					StatusMsg:  err.Error(),
				},
			})
	} else {
		c.JSON(
			http.StatusOK,
			UserLoginResponse{
				Response: Response{
					StatusCode: 0,
				},
				UserId: userid,
				Token:  token,
			})
	}
}

func UserInfo(c *gin.Context) {
	uid := c.Query("user_id")
	token := c.Query("token")

	user, err := service.Info(uid, token)

	// TODO: 确定err
	if user == nil && err != nil {
		c.JSON(
			http.StatusOK,
			UserResponse{
				Response: Response{
					StatusCode: 1,
					StatusMsg:  err.Error(),
				},
			})
	} else {
		c.JSON(
			http.StatusOK,
			UserResponse{
				Response: Response{
					StatusCode: 0,
				},
				User: user,
			})
	}

}
