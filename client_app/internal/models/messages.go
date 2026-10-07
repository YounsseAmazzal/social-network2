package models
type Messagespg struct {
	Last_time string `json:"last_time"` 
	Limit  int `json:"limit"`       
}

func (message *Messagespg) SetMessageDef() {
	if message.Limit <= 0 {
		message.Limit = 10 //  
	}
	if message.Limit > 20 {
		message.Limit = 20 // plaafou
	}
	// the laast time still emty if the user don't send it 

}

type Groupmessages struct{
	ID       string `json:"id,omitempty"`
	GroupID  string `json:"grouupid,omitempty"`
	SenderID string `json:"sender_id,omitempty"`
	GRoup_Message  string `json:"group_message,omitempty"`
	UserAvatar   string `json:"user_avatar,omitempty"`
	FullName string `json:"userFullname,omitempty"`
	Send_date   string `json:"send_date,omitempty"`
}

