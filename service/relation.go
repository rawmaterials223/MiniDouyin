package service

import (
	"fmt"

	"github.com/rawmaterials223/MiniDouyin/repository"
)

// 关注类型全局变量
var DoActionType int = 1   // 关注
var UndoActionType int = 2 // 取消关注

type RelationError struct {
	Status  int
	Message string
}

func (r *RelationError) Error() string {
	return fmt.Sprintf("%v", r.Message)
}

type RelationFlow struct {
	from_user_token string
	to_user_id      int64
	action          int

	from_user_id int64
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
		return &RelationError{
			Status:  1,
			Message: "User doesn't exist",
		} // TODO: 确定err
	}

	f.from_user_id = uid

	//2. 如果action=2，则检查relation表中是否有关注记录
	//3. 如果action=1，relation表中新增关注；如果action=2，删除记录
	relation := &repository.UserRelation{
		FromUserId: f.from_user_id,
		ToUserId:   f.to_user_id,
		IsFollow:   f.action,
	}

	if err := repository.NewRelationDaoInstance().CreateAction(relation); err != nil {
		return err
	}

	fmt.Printf("CreateAction done")

	return nil
}

// 检查token所指用户是否存在
func (f *RelationFlow) IsExistedUser() (bool, int64, error) {

	from_uid, err := repository.NewRelationDaoInstance().QueryUserByToken(f.from_user_token)
	// 用户不存在或查找出错
	if err != nil {
		return false, 0, err
	}

	fmt.Println("QueryUserByToken success")

	return true, from_uid, nil
}
