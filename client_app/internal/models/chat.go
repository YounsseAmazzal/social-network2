package models
type Userchat struct {
	ID              string `json:"id"`
	Nickname        string `json:"nickname"`
	Avatar_URL       string `json:"avatar_url"`
	UserConnected     bool   `json:"User_Connected"`
	LastChatMessage     string `json:"last_chat_message"`
	LastChatDate string `json:"last_message_date"`
	IsOnline        bool   `json:"is_online"`
}