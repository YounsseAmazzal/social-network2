package models

import "time"

type Notification struct {
	ID          string    `json:"id"`
	RecipientID string    `json:"recipient_id"`
	SenderID    string    `json:"sender_id"`
	Category    string    `json:"category"`
	TargetID    string    `json:"target_id"`
	TargetType  string    `json:"target_type"`
	Message     string    `json:"message"`
	IsRead      bool      `json:"is_read"`
	CreatedAt   time.Time `json:"created_at"`
}
type CountermessageNotread struct {
    SenderID string `json:"sender_id"`
    Count    int    `json:"unread_total"`
}
type CountermessageGroupNotread struct {
    GroupID string `json:"group_id"`
    Count   int    `json:"unread_total"`
}