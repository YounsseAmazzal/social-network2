package models

import "time"

type Group struct {
	ID           string    `json:"id"`
	CreatorID    string    `json:"creator_id"` 
	Description  string    `json:"description"`
	Title  string     `json:"group_title"`
	CreatedAt    time.Time `json:"created_at"`
	MembersCount int       `json:"members_count"`
	IsMember     bool      `json:"is_member,omitempty"` 
}
