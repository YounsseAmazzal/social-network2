package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"social/internal/models"
	"social/internal/requests"
	"social/internal/utils"
)

func (handler *Handlers) CreateGroup(res http.ResponseWriter, req *http.Request) {
	//fmt.Println("the user want to CreateGroup ")
	if !requests.IsRequestValid(res, req, "/CreateGroup", http.MethodPost) {
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

	var createGroup models.Group
	if err := json.NewDecoder(req.Body).Decode(&createGroup); err != nil {
		utils.ErrorResponse(res, "JSON NOT VALID", http.StatusBadRequest)
		return
	}
	err = requests.ValidGroupcreation(&createGroup)
	if err != nil {
		utils.ErrorResponse(res, err.Error(), http.StatusBadRequest)
		return
	}
	createGroup.CreatorID = user.ID
	err = handler.services.Reepo.Creategroup(&createGroup)
	if err != nil {
		utils.ErrorResponse(res, err.Error(), http.StatusBadRequest)
		return
	}
	utils.ResponnseSender(res, http.StatusOK, "THe group created succesfully")
	//fmt.Println("that is the user payload", usersingup)

}
func (handler *Handlers) DiscooverGroups(res http.ResponseWriter, req *http.Request) {
	//fmt.Println("the user want to CreateGroup ")
	if !requests.IsRequestValid(res, req, "/DiscoverGroups", http.MethodGet) {
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
	grouups, err := handler.services.Reepo.BringDiscovers(user.ID)
	fmt.Println("errrrrrrrrrrrrr", err)
	if err != nil {
		utils.ErrorResponse(res, "Problem on getting Groups", http.StatusInternalServerError)
		return
	}
	utils.ResponnseSender(res, http.StatusOK, map[string]any{
		"status":        "success",
		"Discovergroup": grouups,
	})

}
func (handler *Handlers) JoinGroup(res http.ResponseWriter, req *http.Request) {
}
func (handler *Handlers) GroupPosts(res http.ResponseWriter, req *http.Request) {
}
func (handler *Handlers) GroupMembers(res http.ResponseWriter, req *http.Request) {
}
func (handler *Handlers) GroupCreateEvent(res http.ResponseWriter, req *http.Request) {
}
func (handler *Handlers) GroupEvents(res http.ResponseWriter, req *http.Request) {
}
func (handler *Handlers) GroupEventParrticipate(res http.ResponseWriter, req *http.Request) {
}
func (handler *Handlers) GroupInvite(res http.ResponseWriter, req *http.Request) {
}
