package models
type PostCreation struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CreaterID    string `json:"createrId"`
	Image       string `json:"shared_image"`
	Date  string `json:"createDate"`
	Postprivecy string `json:"postprivecy"`
	SharedWith []string `json:"only_with":`
}
type Post struct {
	ID            string `json:"id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	CreaterID      string `json:"createrId"`
	AuthorName    string `json:"createrName"`
	AuthorAvatar  string `json:"authorAvatar"`
	CreateDate    string `json:"createDate"`
	Likes         int    `json:"likesNum"`
	Dislikes      int    `json:"dislikesNum"`
	VoteStatus    *int   `json:"vote_status"`
	Image_URL         string `json:"image_url"`
	SharedWith  []int  `json:"omitempty"` 
	Author        User   `json:"creater"`
	TotalComments int    `json:"Commentsnum"`
	Privecytype string `json:"postprivecy"`
}
type PostDetails struct{
	Post
	Raaating *string `json:"userrating"`
	PostComments []*CommentResponse
}
type PostPagination struct {
	Offset int `json:"offset"` // 0000
	Limit  int `json:"limit"`       
}

func (p *PostPagination) SetPostDefaults() {
	if p.Limit <= 0 {
		p.Limit = 5 //  
	}
	if p.Limit > 10 {
		p.Limit = 10 // plaafou
	}
	if p.Offset < 0 {
		p.Offset = 0
	}
}