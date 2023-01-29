package service

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rawmaterials223/MiniDouyin/repository"
	"github.com/rawmaterials223/MiniDouyin/util"
)

func Publish(token, title string, data *multipart.FileHeader) error {
	return NewVideoFlow(token, title, data).Do()
}

func NewVideoFlow(token, title string, data *multipart.FileHeader) *VideoFlow {
	return &VideoFlow{
		token: token,
		title: title,
		data:  data,
	}
}

type Video struct {
	Id            int64  `json:"id"`
	Author        User   `json:"author"`
	PlayUrl       string `json:"play_url"`
	CoverUrl      string `json:"cover_url"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	CommentCount  int64  `json:"comment_count,omitempty"`
	IsFavorite    bool   `json:"is_favorite,omitempty"`
}
type VideoFlow struct {
	token string
	title string
	data  *multipart.FileHeader
}

func (f *VideoFlow) Do() error {
	// 检查用户是否存在
	exist, uid, _ := f.CheckUserByToken()
	if !exist {
		return &ResponseError{1, "User doesn't exist"}
	}

	// 视频文件存储在本地
	filename := filepath.Base(f.data.Filename)
	finalName := fmt.Sprintf("%d_%s", uid, filename)
	saveFile := filepath.Join("./public/video/", finalName)
	util.Logger.Info("video title: " + f.title + ", file path: " + saveFile)

	var c *gin.Context

	if err := c.SaveUploadedFile(f.data, saveFile); err != nil {
		return &RelationError{1, err.Error()}
	}

	video := &repository.Video{
		UserId:     uid,
		PlayUrl:    saveFile,
		CoverUrl:   saveFile,
		Title:      f.title,
		CreateTime: time.Now(),
	}

	err := repository.NewVideoDaoInstance().CreateVideo(video)
	if err != nil {
		return &ResponseError{1, "Upload failed"}
	}

	return nil
}

func (f *VideoFlow) CheckUserByToken() (bool, int64, error) {
	user, err := repository.NewUserDaoInstance().QueryUserByToken(f.token)
	if err != nil {
		return false, 0, err
	}

	util.Logger.Info("QueryUserByToken success")
	return true, user.Id, nil
}
