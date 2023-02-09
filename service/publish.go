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

func PublishList(token, uid string) (VideoList, error) {
	return NewVideoListFlow(token, uid).Do()
}

func NewVideoListFlow(token, uid string) *VideoListFlow {
	s, _ := strconv.ParseInt(uid, 10, 64)
	return &VideoListFlow{
		token:  token,
		userId: s,
	}
}

type VideoList []Video
type VideoResultListType []repository.VideoResult

type Video struct {
	Id            int64          `json:"id"`
	Author        UserResultType `json:"author"`
	PlayUrl       string         `json:"play_url"`
	CoverUrl      string         `json:"cover_url"`
	FavoriteCount int64          `json:"favorite_count"`
	CommentCount  int64          `json:"comment_count"`
	IsFavorite    bool           `json:"is_favorite"`
	Title         string         `json:"title"`
}
type VideoFlow struct {
	token string
	title string
	data  *multipart.FileHeader
}
type VideoListFlow struct {
	token  string
	userId int64

	UserResult      UserResultType
	VideoResultList VideoResultListType
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

func (f *VideoListFlow) Do() (VideoList, error) {

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
		if err := f.QueryUserRelationById(); err != nil {
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

	var videoList VideoList
	for _, v := range f.VideoResultList {
		newV := Video{
			Id:            v.Id,
			Author:        f.UserResult,
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			FavoriteCount: v.FavoriteCount,
			Title:         v.Title,
		}
		videoList = append(videoList, newV)
	}

	return videoList, nil
}

// 查找用户，存在即确认UserInfo
func (f *VideoListFlow) CheckUserByIdToken() (bool, error) {
	_, err := repository.NewUserDaoInstance().QueryUserByIdToken(f.userId, f.token)

	if err != nil {
		return false, err
	}

	util.Logger.Info("QueryUserByIdToken success")

	return true, nil
}

// 查找用户的所有视频，返回videos数组后重新构造
func (f *VideoListFlow) PackVideos() error {
	videoResults, _ := repository.NewVideoDaoInstance().QueryVideoByUid(f.userId)
	f.VideoResultList = videoResults

	return nil
}

// 查找用户的信息，返回包含关注数和粉丝数的用户信息
func (f *VideoListFlow) QueryUserRelationById() error {
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
