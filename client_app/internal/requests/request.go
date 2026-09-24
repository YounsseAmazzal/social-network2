package requests

import (
	"fmt"
	"log"
	"net/http"
	"social/internal/models"
	"social/internal/utils"
	"strings"
)

// insure that the method and the path is good
func IsRequestValid(w http.ResponseWriter, r *http.Request, targetPath, allowedMethod string) bool {
	if strings.Contains(targetPath, "*") {
		pathSegments := strings.Split(r.URL.Path, "/")
		if len(pathSegments) > 0 {
			lastSegment := pathSegments[len(pathSegments)-1]
			targetPath = targetPath[:len(targetPath)-1] + lastSegment
		}
	}
	if r.URL.Path != targetPath {
		// HandleError(w, http.StatusNotFound, "Page not found")
		utils.ErrorResponse(w,"Page not found",http.StatusNotFound)
		log.Printf("404 Path mismatch")
		return false
	}
	if r.Method != allowedMethod {
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Printf("405 Method Not Allowed")
	       utils.ErrorResponse(w,"Method Not Allowed",http.StatusMethodNotAllowed)
		return false
	}
	return true
}
func ValidSignUpPayload(req *http.Request)(models.UserInfo,error){
	var user models.UserInfo
	contentType := req.Header.Get("Content-Type")
	if strings.HasPrefix(contentType, "multipart/form-data") {
		if err := req.ParseMultipartForm(2>>10); err != nil {
			// utils.HandleError(res, http.StatusBadRequest, "Failed to process uploaded data")
			log.Printf("Error parsing multipart form: %v", err)
			return user,fmt.Errorf("Error parsing multipart form: %v", err)
		}
		user.Nickname = req.FormValue("nick_name")
		user.Firstname = req.FormValue("first_name")
		user.Lastname = req.FormValue("last_name")
		user.DateOfBirth = req.FormValue("date_of_birth")

		user.Gender = req.FormValue("gender")
		user.Email = req.FormValue("email")
		user.Password = req.FormValue("password")
		user.ConfirmPassword = req.FormValue("confirm_password")

		if abouuteme:=strings.TrimSpace(req.FormValue("about_me"));abouuteme!=""{
			user.AboutMe=abouuteme
		}


		if file, header, err := req.FormFile("avatar_url"); err == nil {
			defer file.Close()
			if header != nil || header.Size > 0 {
				avatarName, err := utils.ProcessImageUpload(file, header)
				if err != nil {
					return user, fmt.Errorf("failed to save avatar file: %w", err)
				}
				user.Avatar =avatarName
			}
		}

	}
	fmt.Println("thee user retuerner ",user)
	return user,nil
}
func CheckSignupInput(){

}