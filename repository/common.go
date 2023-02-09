package repository

import "time"

// User

// User DB
type User struct {
	Id         int64     `gorm:"column:id" json:"id"`
	Name       string    `gorm:"column:name" json:"name"`
	Token      string    `gorm:"column:token" json:"token,omitempty"`
	CreateTime time.Time `gorm:"column:create_time" json:"create_time,omitempty"`
}

// Userrealtion DB
type UserRelation struct {
	Id         int64 `gorm:"column:id"`
	FromUserId int64 `gorm:"column:from_user_id"`
	ToUserId   int64 `gorm:"column:to_user_id"`
	IsFollow   int   `gorm:"column:is_follow"`
}

// Userrelation Count 勿删，有用
type UserRelationCount struct {
	FollowCount   int  `json:"follow_count"`
	FollowerCount int  `json:"follower_count"`
	IsFollow      bool `json:"is_follow"`
}

// User Result to Client
type UserResult struct {
	Id            int64  `json:"id"`
	Name          string `json:"name"`
	FollowCount   int    `json:"follow_count"`
	FollowerCount int    `json:"follower_count"`
	IsFollow      bool   `json:"is_follow"`
}

// Video

// Video DB
type Video struct {
	Id         int64     `gorm:"column:id"`
	UserId     int64     `gorm:"column:user_id"`
	PlayUrl    string    `gorm:"column:play_url"`
	CoverUrl   string    `gorm:"column:cover_url"`
	Title      string    `gorm:"column:title"`
	CreateTime time.Time `gorm:"column:create_time"`
}

// Videorelation DB
type VideoRelation struct {
	Id         int64 `gorm:"column:id"`
	FromUserId int64 `gorm:"column:from_user_id"`
	ToVideoId  int64 `gorm:"column:to_video_id"`
	IsLike     int   `gorm:"column:is_like"`
}

// Video Result to Client
type VideoResult struct {
	Id int64 `json:"id"`
	//Author        UserResult `json:"author,omitempty"`
	PlayUrl       string `json:"play_url"`
	CoverUrl      string `json:"cover_url"`
	FavoriteCount int64  `json:"favorite_count"`
	CommentCount  int64  `json:"comment_count"`
	IsFavorite    bool   `json:"is_favorite"`
	Title         string `json:"title"`
}
