package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rawmaterials223/MiniDouyin/service"
)

func FavoriteAction(c *gin.Context) {
	token := c.Query("token")
	video_id := c.Query("video_id")
	action_type := c.Query("action_type")

	err := service.FavoriteAction(token, video_id, action_type)
	if err != nil {
		c.JSON(
			http.StatusOK,
			Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	}
}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		//VideoList: DemoVideos,
	})
}
