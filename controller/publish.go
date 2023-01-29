package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rawmaterials223/MiniDouyin/service"
)

type VideoListResponse struct {
	Response
	VideoList []service.Video `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {

	data, err := c.FormFile("data")

	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	title := c.PostForm("title")
	token := c.PostForm("token")

	err = service.Publish(token, title, data)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	} else {
		c.JSON(
			http.StatusOK,
			Response{
				StatusCode: 0,
				StatusMsg:  title + " upload successfully",
			})
	}
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	token := c.Query("token")
	uid := c.Query("user_id")

	videoList, err := service.PublishList(token, uid)
	if err != nil {
		c.JSON(
			http.StatusOK,
			Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			})
	} else {
		c.JSON(
			http.StatusOK,
			VideoListResponse{
				Response: Response{
					StatusCode: 0,
					StatusMsg:  "list",
				},
				VideoList: videoList,
			})
	}

}
