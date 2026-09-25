package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"social/internal/utils"
	"strings"
)

func (handler *Handlers) UserProfile(res http.ResponseWriter, req *http.Request) {
	// fmt.Println("The user want the profile")
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
		utils.ErrorResponse(res, "Your Session Not validd", http.StatusUnauthorized)
		return
	}
	// fmt.Println("********the profilee id**** ",user.ID)
	ProfileeID := extractprofileID(req, user.ID)
	useerProfile, err := handler.services.Reepo.TakeUserProfile(user.ID,ProfileeID)
	if err!=nil{
		if err==sql.ErrNoRows{
					utils.ErrorResponse(res, "Profile nottt found", http.StatusNotFound)
			return 
		}
		utils.ErrorResponse(res, "We Couldn't provide the profile", http.StatusInternalServerError)
	}
	utils.ResponnseSender(res, http.StatusOK, map[string]any{
		"Seeprofile":     useerProfile,
	}) 

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
	// fmt.Println("555555555555555555555555555555555555555")
	UrlPath := strings.TrimPrefix(req.URL.Path, "/UserProfile")
	tackenid1 := strings.Trim(UrlPath, "/")
	if tackenid1 != "" {
		return tackenid1
	}
	// fmt.Println("666666666666666666666666",profileeeid)
	return profileeeid
}
