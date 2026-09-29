package models

type Comment struct {
	ID        string `json:"id"`
	PostID    string `json:"post_id"`
	CreaterID  string `json:"creater_id"`
	Content   string `json:"content"` 
	Image     string `json:"image,omitempty"`
	CreatedAt string `json:"created_at"`
}

type CommentResponse struct {
	ID            string  `json:"id"`
	PostID        string  `json:"post_id"`
	Content       string  `json:"content"`
	CreaterID      string  `json:"creater_id"`
	CreaterName    string  `json:"creater_name"`
	Createravatar  string  `json:"creater_avatar"`
	Image         string  `json:"image_url,omitempty"`
	LikesNum    int     `json:"likes_Num"`
	DislikesNum int     `json:"dislikesNum"`
	UserRate     *string `json:"user_rating"` 
	CreatedAt     string  `json:"created_date"`
}