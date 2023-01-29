package service

import (
	"github.com/rawmaterials223/MiniDouyin/repository"
	"github.com/rawmaterials223/MiniDouyin/util"
)

// 关注类型全局变量
var NoneActionType int = 0 // 无记录
var DoActionType int = 1   // 关注
var UndoActionType int = 2 // 取消关注

type RelationFlow struct {
	from_user_token string
	to_user_id      int64
	action          int

	from_user_id int64
	relation     *repository.UserRelation
}

func RelationAction(from_token string, to_id int64, action int) error {
	return NewRelationFlow(from_token, to_id, action).DoAction()
}

func NewRelationFlow(from_token string, to_id int64, action int) *RelationFlow {
	return &RelationFlow{
		from_user_token: from_token,
		to_user_id:      to_id,
		action:          action,
	}
}

func (f *RelationFlow) DoAction() error {

	//1. 检查用户是否存在
	exist, uid, _ := f.IsExistedUser()
	if !exist {
		return &ResponseError{1, "User doesn't exist"} // TODO: 确定err
	}

	f.from_user_id = uid

	f.relation = &repository.UserRelation{
		FromUserId: f.from_user_id,
		ToUserId:   f.to_user_id,
		IsFollow:   f.action,
	}

	// 查询是否有关注/取消关注记录
	// 不存在/查询结果无，返回0
	// 存在，返回原记录
	action, _ := f.CheckRelationHistory()
	// 确认存在记录，且action不同
	if action != NoneActionType {
		if action != f.action {
			f.UpdateRelation()
			return nil
		} else {
			return &ResponseError{1, "Identical Operation"}
		}
	}

	// 首次关注，添加记录
	if err := f.CreateRelation(); err != nil {
		return err
	}

	util.Logger.Info("RelationAction done")

	return nil
}

// 检查token所指用户是否存在
func (f *RelationFlow) IsExistedUser() (bool, int64, error) {

	from_uid, err := repository.NewRelationDaoInstance().QueryUserByToken(f.from_user_token)
	// 用户不存在或查找出错
	if err != nil {
		return false, 0, err
	}

	return true, from_uid, nil
}

// 查找是否有关注/取消关注记录
func (f *RelationFlow) CheckRelationHistory() (int, error) {
	return repository.NewRelationDaoInstance().QueryRelation(f.from_user_id, f.to_user_id)
}

// 创建新的关注记录
func (f *RelationFlow) CreateRelation() error {
	err := repository.NewRelationDaoInstance().CreateRelation(f.relation)
	if err != nil {
		return err
	}

	return nil
}

// 更新关注/取消关注记录
func (f *RelationFlow) UpdateRelation() error {
	repository.NewRelationDaoInstance().UpdateRelation(f.relation)
	return nil
}
