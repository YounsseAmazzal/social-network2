package models
type ProfileResponse struct {
	User           User           `json:"user"`
	IsPrivate      bool           `json:"isPrivate"`
	IsOwner        bool           `json:"isOwner"`     
	IsFollowing    bool           `json:"isFollowing"`   
	FollowersCount int            `json:"followersCount"`
	FollowingCount int            `json:"followingCount"`
	PostsCount     int            `json:"postsCount"`
	Posts          []PostsOfProfiles  `json:"posts,omitempty"`
	CanSeeTheprofile bool    	`json:"canSee"`
	FollowStatus string `json:"followStatus"`
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
type PostsOfProfiles struct {
	ID            string `json:"id"`
	Title         string `json:"title"`
	Description   string `json:"description"` 
	LikesCount    int    `json:"likes_count"`
	CommentsCount int    `json:"comments_count"`
	CreatedAt     string `json:"created_at"`  
	ImageURL      string `json:"image_url,omitempty"` 
	Privacy       string `json:"privacy,omitempty"`   
}