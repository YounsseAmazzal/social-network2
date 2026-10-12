package requests

import (
	"fmt"
	"log"
	"net/http"
	"social/internal/models"
	"social/internal/utils"
	"strings"
)

func ValidSignUpPayload(req *http.Request) (models.UserInfo, error) {
	var user models.UserInfo
	contentType := req.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "multipart/form-data") {
		if err := req.ParseMultipartForm(32 << 20); err != nil {
			// utils.HandleError(res, http.StatusBadRequest, "Failed to process uploaded data")
			log.Printf("Error parsing multipart form: %v", err)
			return user, fmt.Errorf("Error parsing multipart form: %v", err)
		}

		user.Nickname = req.FormValue("nick_name")
		user.Firstname = req.FormValue("first_name")
		user.Lastname = req.FormValue("last_name")
		user.DateOfBirth = req.FormValue("date_of_birth")

		user.Gender = req.FormValue("gender")
		user.Email = req.FormValue("email")
		user.Password = req.FormValue("password")
		user.ConfirmPassword = req.FormValue("confirm_password")

		if abouuteme := strings.TrimSpace(req.FormValue("about_me")); abouuteme != "" {
			user.AboutMe = abouuteme
		}

		if file, header, err := req.FormFile("avatar"); err == nil {
			defer file.Close()
			if header != nil || header.Size > 0 {
				// log.Fatal("here come herererererer")
				avatarName, err := utils.ProcessImageUpload(file, header,"auth")
				if err != nil {
					return user, fmt.Errorf("failed to save avatar file: %w", err)
				}
				user.Avatar = avatarName
			}
		}

	}
	//fmt.Println("@@@@@@@@thee user retuerner ", user.DateOfBirth)
	return user, nil
}
func CheckSignupInput() {

}
