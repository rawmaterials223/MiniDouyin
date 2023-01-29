package service

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"sync"
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

func PublishList(token, uid string) ([]Video, error) {
	return NewVideoListFlow(token, uid).Do()
}

func NewVideoListFlow(token, uid string) *VideoListFlow {
	s, _ := strconv.ParseInt(uid, 10, 64)
	return &VideoListFlow{
		token:  token,
		userId: s,
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
	Title         string `json:"title,omitempty"`
}
type VideoFlow struct {
	token string
	title string
	data  *multipart.FileHeader
}
type VideoListFlow struct {
	token  string
	userId int64

	UserInfo     *repository.User
	UserRelation *repository.UserRelationCount
	VideoList    []repository.Video
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
	// 尝试并行
	var wg sync.WaitGroup
	wg.Add(2)
	var saveFileErr, createErr error
	go func() {
		defer wg.Done()
		var c *gin.Context
		if err := c.SaveUploadedFile(f.data, saveFile); err != nil {
			saveFileErr = &ResponseError{1, err.Error()}
			return
		}
	}()
	go func() {
		defer wg.Done()
		video := &repository.Video{
			UserId:     uid,
			PlayUrl:    saveFile,
			CoverUrl:   saveFile, // TODO: CoverUrl待定
			Title:      f.title,
			CreateTime: time.Now(),
		}

		if err := repository.NewVideoDaoInstance().CreateVideo(video); err != nil {
			createErr = &ResponseError{1, f.title + " upload failed"}
			return
		}
	}()
	wg.Wait()
	if saveFileErr != nil {
		return saveFileErr
	} else if createErr != nil {
		return createErr
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

func (f *VideoListFlow) Do() ([]Video, error) {

	exist, _ := f.CheckUserByIdToken()
	if !exist {
		return nil, &ResponseError{1, "User doesn't exist"}
	}

	//尝试并行
	var wg sync.WaitGroup
	wg.Add(2)
	var authorErr, videoErr error
	// Author 信息
	go func() {
		defer wg.Done()
		if err := f.PackRelation(); err != nil {
			authorErr = err
			return
		}
	}()
	// Video信息
	go func() {
		defer wg.Done()
		if err := f.PackVideos(); err != nil {
			videoErr = err
			return
		}

	}()
	wg.Wait()
	if authorErr != nil {
		return nil, &ResponseError{1, "author list error"}
	} else if videoErr != nil {
		return nil, &ResponseError{1, "video list error"}
	}

	author := User{
		Id:            f.userId,
		Name:          f.UserInfo.Name,
		FollowCount:   f.UserRelation.FollowCount,
		FollowerCount: f.UserRelation.FollowerCount,
		IsFollow:      f.UserRelation.IsFollow,
	}

	var VideoList []Video
	for _, v := range f.VideoList {
		newV := Video{
			Id:       v.Id,
			Author:   author,
			PlayUrl:  v.PlayUrl,
			CoverUrl: v.CoverUrl,
			Title:    v.Title,
		}
		VideoList = append(VideoList, newV)
	}

	return VideoList, nil
}

// 查找用户，存在即确认UserInfo
func (f *VideoListFlow) CheckUserByIdToken() (bool, error) {
	user, err := repository.NewUserDaoInstance().QueryUserByIdToken(f.userId, f.token)

	if err != nil {
		return false, err
	}

	util.Logger.Info("QueryUserByIdToken success")
	f.UserInfo = user

	return true, nil
}

// 查找用户信息，确认UserRelation
func (f *VideoListFlow) PackRelation() error {
	follow_count, follower_count, _ := repository.NewRelationDaoInstance().CalculateRelation(f.userId)
	f.UserRelation = &repository.UserRelationCount{
		FollowCount:   int(follow_count),
		FollowerCount: int(follower_count),
		IsFollow:      true,
	}

	return nil
}

// 查找用户的所有视频，返回videos数组后重新构造
func (f *VideoListFlow) PackVideos() error {
	videos, _ := repository.NewVideoDaoInstance().QueryVideoByUid(f.userId)
	f.VideoList = videos

	return nil
}
