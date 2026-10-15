package models

import (
	"time"
)

type GroupItem struct {
	ID           string    `json:"id"`
	CreatorID    string    `json:"creator_id"`
	Description  string    `json:"description"`
	Title        string    `json:"group_title"`
	CreatedAt    time.Time `json:"created_at"`
	MembersCount int       `json:"members_count"`
	IsMember     bool      `json:"is_member,omitempty"`
	IsRequested  string      `json:"request_status,omitempty"`
	GroupMember  []string  `json:"group_members,omitempty"`
}
type GroupDetail struct {
	ID           string    `json:"id"`
	CreatorID    string    `json:"creator_id"`
	Description  string    `json:"description"`
	Title        string    `json:"group_title"`
	CreatedAt    time.Time `json:"created_at"`
	MembersCount int       `json:"members_count"`
	IsMember     bool      `json:"is_member,omitempty"`
	Creater      User      `json:"creater"`
}
type GroupPostCreation struct {
	GrouupID  string `json:"group_id"`
	GroupPost PostCreation
}
type GroupCommentCreation struct {
	GrouupID  string `json:"group_id"`
	GroupPost PostCreation
}
type GrouupPR struct {
	Offset int `json:"offset"` // 0000
	Limit  int `json:"limit"`
}

func (p *GrouupPR) SetDefaults() {
	if p.Limit <= 0 {
		p.Limit = 10 //
	}
	if p.Limit > 100 {
		p.Limit = 100 // plaafou
	}
	if p.Offset < 0 {
		p.Offset = 0
	}
}

type GrouupPostDetails struct {
	Post
	Raaating     *string `json:"userrating"`
	PostComments []*GroupCommentResponse
}
type GroupCommentResponse struct {
	ID            string `json:"id"`
	PostID        string `json:"post_id"`
	Content       string `json:"content"`
	CreaterID     string `json:"creater_id"`
	CreaterName   string `json:"creater_name"`
	Createravatar string `json:"creater_avatar"`
	Image         string `json:"image_url,omitempty"`
	CreatedAt     string `json:"created_date"`
}

type GroupComment struct {
	ID        string `json:"id"`
	PostID    string `json:"post_id"`
	CreaterID string `json:"creater_id"`
	Content   string `json:"content"`
	Image     string `json:"image,omitempty"`
	CreatedAt string `json:"created_at"`
}
