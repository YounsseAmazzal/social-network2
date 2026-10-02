package models

import "time"

type Eventplyload struct {
	GrouupID     string     `json:"group_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	EventDate   time.Time `json:"event_date"`
}
type GroupEvent struct {
	ID            string    `json:"id"`
	GroupID       string    `json:"group_id"`
	CreatorID     string    `json:"creator_id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	EventDate     time.Time `json:"event_date"`  
	CreatedAt     time.Time `json:"created_at"`   
	userrate      string    `json:"user_rate,omitempty"` 
	TotalGoing    int       `json:"total_participate"`
	TotalNotGoing int       `json:"total_participate"`
	Creator        User      `json:"Creator,omitempty"`
}