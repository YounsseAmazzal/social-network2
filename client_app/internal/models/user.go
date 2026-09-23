package models
type UserInfo struct {
	ID              string `json:"id"`
	Firstname       string `json:"first_name"`
	Lastname        string `json:"last_name"`
	Nickname        string `json:"nick_name"`
	DateOfBirth     string `json:"date_of_birth"`
	Age             int    `json:"age"`
	Gender          string `json:"gender"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	AvatarURL       string `json:"avatar_url,omitempty"`
	Avatar          string `json:"avatar,omitempty"`
	AboutMe         string `json:"about_me"`
	IsPrivate       int    `json:"is_Private"`
}
type User struct {
	ID         string `json:"id"`
	Nickname   string `json:"nickname"`
	Username   string `json:"username,omitempty"`
	Firstname  string `json:"firstname"`
	Lastname   string `json:"lastname"`
	Age        int    `json:"age"`
	Gender     string `json:"gender"`
	IsLoggedIn bool   `json:"is_logged_in"`
	Password        string `json:"-"`
	Email      string `json:"email"`
	AvatarURL  string `json:"avatar_url"`
	Avatar     any    `json:"avatar,omitempty"`
	AboutMe    string `json:"aboutMe,omitempty"`
	IsPrivate  bool   `json:"isPrivate"`
}
type UserLoign struct {
	Identifiant string
	Password    string
}