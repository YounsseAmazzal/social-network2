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
	Userrate      string    `json:"user_rate,omitempty"` 
	TotalGoing    int       `json:"total_participate"`
	TotalNotGoing int       `json:"total_Not_participate"`
	Creator        User      `json:"Creator,omitempty"`
}

type EventPagination struct {
	Offset int `json:"offset"` // 0000
	Limit  int `json:"limit"`       
}

func (p *EventPagination) SetEventDefaults() {
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