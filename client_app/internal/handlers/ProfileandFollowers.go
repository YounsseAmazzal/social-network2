package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"social/internal/models"
	"social/internal/requests"
	"social/internal/utils"
	"strings"
)

func (handler *Handlers) UserProfile(res http.ResponseWriter, req *http.Request) {
	// fmt.Println("The user want the profile")

	if !requests.IsRequestValid(res, req, "/UserProfile", http.MethodGet) {
		return
	}
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
	useerProfile, err := handler.services.Reepo.TakeUserProfile(user.ID, ProfileeID)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.ErrorResponse(res, "Profile nottt found", http.StatusNotFound)
			return
		}
		utils.ErrorResponse(res, "We Couldn't provide the profile", http.StatusInternalServerError)
	}
	utils.ResponnseSender(res, http.StatusOK, map[string]any{
		"Seeprofile": useerProfile,
	})

}
func (handler *Handlers) ProfileUppdate(res http.ResponseWriter, req *http.Request) {
	fmt.Println("The user want to Update the profile")
	if !requests.IsRequestValid(res, req, "/ProfileUppdate", http.MethodPatch) {
		return
	}
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
	var payload models.UpdateProfileRequest
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		utils.ErrorResponse(res, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	newNickname := user.Nickname
	newAboutMe := user.AboutMe
	newIsPrivate := user.IsPrivate
	if payload.Nickname != nil {
		cleanNick := strings.ToLower(strings.TrimSpace(*payload.Nickname))
		if cleanNick != "" {
			newNickname = html.EscapeString(cleanNick)
		}
	}
	if payload.AboutMe != nil {
		cleanAbout := strings.TrimSpace(*payload.AboutMe)
		if len(cleanAbout) > 200 {
			utils.ErrorResponse(res, "About me too long (max 200)", http.StatusBadRequest)
			return
		}
		newAboutMe = html.EscapeString(cleanAbout)
	}
	if payload.IsPrivate != nil {
		newIsPrivate = *payload.IsPrivate
	}
	err = handler.services.Reepo.UpdateProfileSettings(user.ID, newNickname, newAboutMe, newIsPrivate)
	if err != nil {
		// (Unique Constraint)
		if strings.Contains(err.Error(), "UNIQUE") {
			utils.ErrorResponse(res, "Nickname already taken", http.StatusConflict)
			return
		}
		utils.ErrorResponse(res, "Database update failed", http.StatusInternalServerError)
		return
	}
	utils.ResponnseSender(res, http.StatusOK, map[string]any{
		"message": "Profile updated successfully",
		"user": map[string]any{
			"nickname":   newNickname,
			"about_me":   newAboutMe,
			"is_private": newIsPrivate,
		},
	})

}
func (handler *Handlers) FollowersAccept(res http.ResponseWriter, req *http.Request) {
	fmt.Println("The user want the FollowersAccept")
	if !requests.IsRequestValid(res, req, "/followers/Accepting", http.MethodPost) {
		return
	}
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
	followedID, err := requests.Accepting(req, user.ID)
	if err != nil {
		utils.ErrorResponse(res, err.Error(), http.StatusBadRequest)
		return
	}
	// fmt.Println("Heee came here ")
	pendding, errr := handler.services.Reepo.CheckPending(user.ID, followedID)
	if errr != nil {
		utils.ErrorResponse(res, "Reequeest invalid", http.StatusInternalServerError)
		return
	}
	if !pendding {
		utils.ErrorResponse(res, "THere is no pendding requests heere", http.StatusInternalServerError)
		return
	}
	err = handler.services.Reepo.AcceptFollowRequest(user.ID, followedID)
	if err != nil {
		utils.ErrorResponse(res, err.Error(), http.StatusInternalServerError)
		return
	}
	nicknamenothif := user.ID
	if user.Nickname != "" {
		nicknamenothif = user.Nickname
	}
	fmt.Println("the name of sender", nicknamenothif)

	///i will add the notifications here
	utils.ResponnseSender(res, http.StatusOK, "Follow request accepted")
}
func (handler *Handlers) FollowersPendiing(res http.ResponseWriter, req *http.Request) {
	fmt.Println("The user want the profile")
}
func (handler *Handlers) FollowersDeclineinng(res http.ResponseWriter, req *http.Request) {
	fmt.Println("The user want FollowersDeclineinng")
	if !requests.IsRequestValid(res, req, "/followers/Declineinng", http.MethodPost) {
		return
	}
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
	followedID, err := requests.Declineinng(req, user.ID)
	if err != nil {
		utils.ErrorResponse(res, err.Error(), http.StatusBadRequest)
		return
	}
	// fmt.Println("Heee came here ")
	pendding, errr := handler.services.Reepo.CheckPending(followedID, user.ID)
	if errr != nil {
		utils.ErrorResponse(res, "Reequeest invalid", http.StatusInternalServerError)
		return
	}
	fmt.Println("The follow pendding",pendding)
	if !pendding {
		utils.ErrorResponse(res, "THere is no pendding requests heere", http.StatusInternalServerError)
		return
	}
	err = handler.services.Reepo.Unfollowing(user.ID, followedID)
	if err != nil {
		utils.ErrorResponse(res, "The request of decline refoced no one to decline", http.StatusInternalServerError)
		return
	}
	message := "Follow Request  is Declined"
	utils.ResponnseSender(res, http.StatusOK, map[string]any{
		"status":  "success",
		"message": message,
	})
}
func (handler *Handlers) UserUnfollow(res http.ResponseWriter, req *http.Request) {
	fmt.Println("The user want UserUnfollow")
	if !requests.IsRequestValid(res, req, "/UserUnfollow", http.MethodPost) {
		return
	}
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
	followedID, err := requests.ValidUnfollow(req, user.ID)
	if err != nil {
		utils.ErrorResponse(res, err.Error(), http.StatusBadRequest)
		return
	}
	ProfileTarget, err := handler.services.Reepo.TakeUserProfile(user.ID, followedID)
	// _, err = handler.services.Reepo.TakeUserProfile(user.ID, followedID)
	if err != nil {
		utils.ErrorResponse(res, err.Error(), http.StatusUnauthorized)
		return
	}
	ProfileTarget.IsPrivate = true
	message := "You unFollow the user"
	err = handler.services.Reepo.Unfollowing(user.ID, followedID)
	if err != nil {
		utils.ErrorResponse(res, err.Error(), http.StatusUnauthorized)
		return
	}
	utils.ResponnseSender(res, http.StatusOK, map[string]any{
		"status":  "success",
		"message": message,
	})

}
func (handler *Handlers) Userfollowinng(res http.ResponseWriter, req *http.Request) {
	if !requests.IsRequestValid(res, req, "/Userfollowinng", http.MethodPost) {
		return
	}
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

	// TheFollower:=user.ID
	// TheFollowed:=req.URL.Query().Get("Userfollowinng_id")
	// if TheFollowed==""{
	// 	utils.ErrorResponse(res, "The taget follow not found", http.StatusUnauthorized)
	// 		return
	// }
	// if TheFollowed==TheFollower{
	// 	utils.ErrorResponse(res, "You Couldn't follow your self hh", http.StatusBadRequest)
	// 	return
	// }
	followedID, err := requests.ValidFollowReq(req, user.ID)
	if err != nil {
		utils.ErrorResponse(res, err.Error(), http.StatusBadRequest)
		return
	}
	// fmt.Println("HEEEEEEEEEEEEEEREEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEe",followedID)
	ProfileTarget, err := handler.services.Reepo.TakeUserProfile(user.ID, followedID)
	fmt.Println("$$$$$$$$$$$$$$$$Check the Privicy$$$$$$$", ProfileTarget.IsPrivate)
	ProfileTarget.IsPrivate = true
	// _, err = handler.services.Reepo.TakeUserProfile(user.ID, followedID)
	if err != nil {
		utils.ErrorResponse(res, err.Error(), http.StatusUnauthorized)
		return
	}
	// if !ProfileTarget.IsPrivate{
	// 	//handelt the accepten and the nothification
	// 	fmt.Println("!!!!!!!!!!!!!!!Profile is Privite!!!!!!!!!!!!!")

	// }
	Followstatuss, err := handler.services.Reepo.Follwing(user.ID, followedID)
	if Followstatuss == "followpending" {
		fmt.Println("The follow is Penddding")
		//// push the nothification
		utils.ResponnseSender(res,http.StatusOK, map[string]any{
			"message":"The follow penddding while unitil the user accept"})
		return
	}
	message := "You are now following this user"
	if Followstatuss == "followpending" {
		message = "Follow request sent and pending approval"
	}
	//nothification
	utils.ResponnseSender(res, http.StatusOK, map[string]any{
		"status":  "success",
		"message": message,
		"data": map[string]string{
			"follow_status": Followstatuss,
		},
	})
}
func (handler *Handlers) Userfollowers(res http.ResponseWriter, req *http.Request) {
	fmt.Println("The user want the profile")
}
func extractprofileID(req *http.Request, profileeeid string) string {
	tackenid := strings.TrimSpace(req.URL.Query().Get("id"))
	if tackenid != "" {
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
