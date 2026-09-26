package models
type ProfileResponse struct {
	User           User           `json:"user"`
	IsPrivate      bool           `json:"isPrivate"`
	IsOwner        bool           `json:"isOwner"`     
	IsFollowing    bool           `json:"isFollowing"`   
	FollowersCount int            `json:"followersCount"`
	FollowingCount int            `json:"followingCount"`
	PostsCount     int            `json:"postsCount"`
	Posts          []PostPreview  `json:"posts,omitempty"`
}

type PostPreview struct {
	ID            string `json:"id"`
	Title         string `json:"title"`
	Content       string `json:"content"`
	ImageURL      string `json:"imageUrl,omitempty"`
	LikesCount    int    `json:"likesCount"`
	CommentsCount int    `json:"commentsCount"`
	CreatedAt     string `json:"createdAt"`
}
type UpdateProfileRequest struct {
	Nickname  *string `json:"nickname"`
	AboutMe   *string `json:"about_me"`
	IsPrivate *bool   `json:"is_private"`
}