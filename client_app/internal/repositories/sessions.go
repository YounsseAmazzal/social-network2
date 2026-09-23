package repositories

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
)
const (
	SessionCookieName = "solo_social"
	SessionDuration   = 24 * time.Hour
)

func (sessiorepo *Repository) CreateUserSession(w http.ResponseWriter, userID string) error {
	tokenID, err := uuid.NewV4()
	if err != nil {
		return fmt.Errorf("session: failed to generate token: %w", err)
	}
	token := tokenID.String()
	expiry := time.Now().Add(SessionDuration)

	// i prefer trasation to manage data base n
	_, err = sessiorepo.DB.Database.Exec("DELETE FROM sessions WHERE user_id = ?", userID)
	if err != nil {
		return fmt.Errorf("session: failed to clear old sessions: %w", err)
	}

	query := "INSERT INTO sessions (token, user_id, expire_at) VALUES (?, ?, ?)"
	_, err = sessiorepo.DB.Database.Exec(query, token, userID, expiry)
	if err != nil {
		return fmt.Errorf("session: failed to save session to db: %w", err)
	}

	// (Coooookie ssssetting)
	setAuthCookie(w, token, expiry)

	return nil
}
//session 
func setAuthCookie(w http.ResponseWriter, token string, expiry time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    token,
		Expires:  expiry,
		Path:     "/",
		HttpOnly: true,  ///no xsss 
		SameSite: http.SameSiteLaxMode, // avoiiid CSRF
		Secure:   false,                
	})
}
