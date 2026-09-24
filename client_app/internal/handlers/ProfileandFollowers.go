package handlers

import (
	"fmt"
	"net/http"
	"social/internal/utils"
	"strings"
)

func (handler *Handlers) UserProfile(res http.ResponseWriter, req *http.Request) {
	fmt.Println("The user want the profile")
	cookie, err := req.Cookie("solo_social")
	if err != nil {
		utils.ErrorResponse(res, "sesssion not valiiid ", http.StatusUnauthorized)
		return
	}
	_, err = handler.services.Reepo.Checksession(cookie.Value)
	if err != nil {
		utils.ErrorResponse(res, "Your Session Not valid", http.StatusUnauthorized)
		// log.Printf("Your Session Not valid", err.Error())
		return
	}
	user, err := handler.services.Reepo.GetUserByToken(cookie.Value)
	if err != nil {
		utils.ErrorResponse(res, "Your Session Not valid", http.StatusUnauthorized)
		return
	}
	ProfileeID := extractprofileID(req, user.ID)

	useerProfile, err := handler.services.Reepo.TakeUserProfile(user.ID,ProfileeID)

}
func (handler *Handlers) ProfileUppdate(res http.ResponseWriter, req *http.Request) {
	fmt.Println("The user want the profile")
}
func (handler *Handlers) FollowersAccept(res http.ResponseWriter, req *http.Request) {
	fmt.Println("The user want the profile")
}
func (handler *Handlers) FollowersPendiing(res http.ResponseWriter, req *http.Request) {
	fmt.Println("The user want the profile")
}
func (handler *Handlers) FollowersDeclineinng(res http.ResponseWriter, req *http.Request) {
	fmt.Println("The user want the profile")
}
func (handler *Handlers) UserUnfollow(res http.ResponseWriter, req *http.Request) {
	fmt.Println("The user want the profile")
}
func (handler *Handlers) Userfollowinng(res http.ResponseWriter, req *http.Request) {
	fmt.Println("The user want the profile")
}
func (handler *Handlers) Userfollowers(res http.ResponseWriter, req *http.Request) {
	fmt.Println("The user want the profile")
}
func extractprofileID(req *http.Request, profileeeid string) string {
	tackenid := strings.TrimSpace(req.URL.Query().Get("id"))
	if tackenid!=""{ 
		return tackenid
	}

	UrlPath := strings.TrimPrefix(req.URL.Path, "/profile")
	tackenid1 := strings.Trim(UrlPath, "/")
	if tackenid1 != "" {
		return tackenid1
	}

	return profileeeid
}
