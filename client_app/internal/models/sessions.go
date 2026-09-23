package models

import "time"

type sessions struct {
	Tocken    string    `json:"id"`
	User_ID    string    `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
	CreateedAt time.Time `json:"created_at"`
}
