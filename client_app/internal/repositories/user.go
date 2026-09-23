package repositories

import (
	"fmt"
	"html"
	"social/internal/models"
	"social/internal/utils"
	"strings"

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
