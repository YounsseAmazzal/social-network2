package repositories

import (
	"database/sql"
	"fmt"
	"html"
	"log"
	"social/internal/models"
	"social/internal/utils"
	"strings"
	"time"

	"github.com/gofrs/uuid"
)

func (ur *Repository) CreateUser(payload *models.UserInfo) (*models.User, error) {
	newID, err := uuid.NewV4()
	if err != nil {
		return nil, fmt.Errorf("repository: failed to generate uuid: %w", err)
	}
	user := &models.User{
		ID:        newID.String(),
		Firstname: strings.TrimSpace(payload.Firstname),
		Lastname:  strings.TrimSpace(payload.Lastname),
		Email:     strings.ToLower(strings.TrimSpace(payload.Email)),
		Password:  payload.Password,
		Age:       payload.Age,
		Gender:    payload.Gender,
		AboutMe:   html.EscapeString(payload.AboutMe),
		IsPrivate: payload.IsPrivate == 1,
	}
	nickname := strings.TrimSpace(payload.Nickname)
	if nickname == "" {
		if len(user.Firstname) > 0 {
			user.Nickname = strings.ToLower(string(user.Firstname[0]) + user.Lastname)
		} else {
			user.Nickname = "user_" + user.ID[:8]
		}
	} else {
		user.Nickname = strings.ToLower(nickname)
	}

	user.AvatarURL = payload.AvatarURL
	if user.AvatarURL == "" {
		user.AvatarURL = utils.Defauteavaaatar
	}

	query := `
		INSERT INTO users (
			id, nickname, firstname, lastname, age, 
			gender, email, password, avatarURL, about_me, is_private
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err = ur.DB.Database.Exec(query,
		user.ID,
		user.Nickname,
		user.Firstname,
		user.Lastname,
		user.Age,
		user.Gender,
		user.Email,
		user.Password,
		user.AvatarURL,
		user.AboutMe,
		payload.IsPrivate,
		//autmaaateee privete (boool )
	)
	if err != nil {
		return nil, fmt.Errorf("repository: failed to insert user: %w", err)
	}
	return user, nil
}
func (ur *Repository) CheckuserIfExxist(usernameoremail string) (*models.User, bool) {
	var UserInfo models.User
	fmt.Println("the indefair", usernameoremail)
	usernameoremail = strings.ToLower(usernameoremail)
	query := `
		SELECT 
			id, nickname, firstname, lastname, age, gender, 
			email, avatarURL, COALESCE(about_me, ''), COALESCE(is_private, 0), password 
		FROM users 
		WHERE nickname = ? OR email = ? 
		LIMIT 1`

	err := ur.DB.Database.QueryRow(query, usernameoremail, usernameoremail).Scan(
		&UserInfo.ID,
		&UserInfo.Nickname,
		&UserInfo.Firstname,
		&UserInfo.Lastname,
		&UserInfo.Age,
		&UserInfo.Gender,
		&UserInfo.Email,
		&UserInfo.AvatarURL,
		&UserInfo.AboutMe,
		&UserInfo.IsPrivate,
		&UserInfo.Password,
	)
	fmt.Println("that is the id of user", UserInfo.ID)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, false
		}
		log.Printf("Repository Error Fach bghit njib user:%v", err)
		return nil, false
	}

	return &UserInfo, true
}
func (ur *Repository) Checksession(tocccken string) (string, error) {
	var userID string
	var expireAt time.Time
	//queery to data base 
	query := "SELECT user_id, expire_at FROM sessions WHERE token = ?"
	err := ur.DB.Database.QueryRow(query, tocccken).Scan(&userID, &expireAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("session not found")
		}
		return "", fmt.Errorf("database error: %w", err)
	}

	if time.Now().After(expireAt) {
		//cleaaant the dataa base 
		_ = ur.DeleteSession(tocccken)
		return "", fmt.Errorf("session expired")
	}

	return userID, nil
}
func (ur *Repository) DeleteSession(token string) error {
	query := "DELETE FROM sessions WHERE token = ?"
	_, err := ur.DB.Database.Exec(query, token)
	if err != nil {
		return fmt.Errorf("repository: failed to delete session: %w", err)
	}
	return nil
}