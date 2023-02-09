package repository

// User

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

// Video Result to Client
type VideoResult struct {
	Id            int64  `json:"id"`
	PlayUrl       string `json:"play_url"`
	CoverUrl      string `json:"cover_url"`
	FavoriteCount int64  `json:"favorite_count"`
	CommentCount  int64  `json:"comment_count"`
	IsFavorite    bool   `json:"is_favorite"`
	Title         string `json:"title"`
}
