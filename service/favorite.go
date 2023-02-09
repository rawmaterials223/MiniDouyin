package service

import (
	"strconv"

	"github.com/rawmaterials223/MiniDouyin/repository"
	"github.com/rawmaterials223/MiniDouyin/util"
)

// 赞类型全局变量
var NoneFavType int = 0 // 无记录
var FavType int = 1     // 赞
var UndoFavType int = 2 // 取消赞

func FavoriteAction(token, vid, action string) error {
	return NewFavoriteFlow(token, vid, action).DoAction()
}

func NewFavoriteFlow(token, vid, action string) *FavoriteFlow {
	s, _ := strconv.ParseInt(vid, 10, 64)
	t, _ := strconv.Atoi(action)

	return &FavoriteFlow{
		token:   token,
		videoId: s,
		action:  t,
	}
}

type FavoriteFlow struct {
	token   string
	videoId int64
	action  int

	from_user_id  int64
	videoRelation *repository.VideoRelation
}

func (f *FavoriteFlow) DoAction() error {

	// 检查用户是否存在
	exist, uid, _ := f.CheckUserByToken()
	if !exist {
		return &ResponseError{1, "User doesn't exist"}
	}
	f.from_user_id = uid

	// 检查视频是否存在
	exist, _ = f.CheckVideoById()
	if !exist {
		return &ResponseError{1, "Video doesn't exist"}
	}

	f.videoRelation = &repository.VideoRelation{
		FromUserId: f.from_user_id,
		ToVideoId:  f.videoId,
		IsLike:     f.action,
	}

	// 查询是否有赞记录
	// 不存在，返回0
	// 存在，返回原记录
	action, _ := f.CheckFavoriteHistory()
	util.Logger.Info("action is ")
	if action != NoneFavType {
		if action != f.action {
			f.UpdateFavorite()
			return nil
		} else {
			return &ResponseError{1, "Identical Operation"}
		}
	}

	// 首次赞，添加记录
	if err := f.CreateFavorite(); err != nil {
		return &ResponseError{1, "Create favorite record failed"}
	}

	util.Logger.Info("FavoriteAction Done")
	return nil
}

// 查找用户是否存在
func (f *FavoriteFlow) CheckUserByToken() (bool, int64, error) {
	user, err := repository.NewUserDaoInstance().QueryUserByToken(f.token)
	if err != nil {
		return false, 0, err
	}

	util.Logger.Info("QueryUserByToken success")
	return true, user.Id, nil
}

// 查找视频是否存在
func (f *FavoriteFlow) CheckVideoById() (bool, error) {
	video, _ := repository.NewVideoDaoInstance().QueryVideoById(f.videoId)
	if video == nil {
		return false, nil
	}
	return true, nil
}

// 查找是否有赞操作
func (f *FavoriteFlow) CheckFavoriteHistory() (int, error) {
	return repository.NewVideoRelationDaoInstance().QueryRelation(f.from_user_id, f.videoId)
}

// 更新赞/取消赞操作
func (f *FavoriteFlow) UpdateFavorite() error {
	return repository.NewVideoRelationDaoInstance().UpdateRelation(f.videoRelation)
}

// 创建赞操作
func (f *FavoriteFlow) CreateFavorite() error {
	err := repository.NewVideoRelationDaoInstance().CreateRelation(f.videoRelation)
	if err != nil {
		return err
	}
	return nil
}
