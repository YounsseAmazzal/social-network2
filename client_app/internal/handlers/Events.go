package handlers

import (
	"fmt"
	"log"
	"net/http"
	"social/internal/requests"
	"social/internal/utils"
)

func (handler *Handlers) GroupCreateEvent(res http.ResponseWriter, req *http.Request) {
	if !requests.IsRequestValid(res, req, "/GroupCreateEvent", http.MethodPost) {
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
	Eventplyload, err := requests.Eventplyload(req)
	err = requests.Vaalidevent(&Eventplyload)
	if err != nil {
		//fmt.Println("YAAAAAAAAAAAK", err)
		utils.ErrorResponse(res, err.Error(), http.StatusUnauthorized)
		return
	}
	grouupid, err := requests.ExtraactGroupId(req)
	if err != nil {
		utils.ErrorResponse(res, "Grouuup not Found", http.StatusUnauthorized)
		return
	}
	Groupexist, err := handler.services.Reepo.GoupExist(grouupid)
	if !Groupexist {
		utils.ErrorResponse(res, "Grouup Not Fouund", http.StatusNotFound)
		return
	}
	ismember, errr := handler.services.Reepo.CheckifUsserMember(user.ID, grouupid)
	fmt.Println("The user is Member ?????", ismember)
	if errr != nil {
		fmt.Println("THE MEEEEMER CHECB err", errr)
		utils.ErrorResponse(res, "Problem on checking membership", http.StatusNotFound)
		return
	}
	if !ismember {
		utils.ErrorResponse(res, "You  are not a member to see this group", http.StatusNotFound)
		return
	}
	Eventplyload.GrouupID=grouupid
	eventCreated,err:= handler.services.Reepo.CreateEvent(&Eventplyload,user.ID)
	if err != nil {
		log.Fatal("hi")
		utils.ErrorResponse(res, err.Error(), http.StatusInternalServerError)
		// utils.ErrorResponse(res, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.ResponnseSender(res, http.StatusOK, map[string]any{
		"message": "The Group post created with success",
		"eventCreated":eventCreated,
		"status":  "success",
	})

}
func (handler *Handlers) GroupEvents(res http.ResponseWriter, req *http.Request) {
	if !requests.IsRequestValid(res, req, "/GroupEvents", http.MethodGet) {
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
	grouupid, err := requests.ExtraactGroupId(req)
	if err != nil {
		utils.ErrorResponse(res, "Grouuup not Found", http.StatusUnauthorized)
		return
	}
	Groupexist, err := handler.services.Reepo.GoupExist(grouupid)
	if !Groupexist {
		utils.ErrorResponse(res, "Grouup Not Fouund", http.StatusNotFound)
		return
	}
	ismember, errr := handler.services.Reepo.CheckifUsserMember(user.ID, grouupid)
	fmt.Println("The user is Member ?????", ismember)
	if errr != nil {
		fmt.Println("THE MEEEEMER CHECB err", errr)
		utils.ErrorResponse(res, "Problem on checking membership", http.StatusNotFound)
		return
	}
	if !ismember {
		utils.ErrorResponse(res, "You  are not a member to see this group", http.StatusNotFound)
		return
	}
	GetGroupEvents,err:=handler.services.Reepo.GetGroupEvents(grouupid,user.ID)
		if err != nil {
		utils.ErrorResponse(res, err.Error(), http.StatusInternalServerError)
		// utils.ErrorResponse(res, err.Error(), http.StatusInternalServerError)
		return
	}
	// fmt.Println("Those is Grouuup Events",GetGroupEvents)
	utils.ResponnseSender(res, http.StatusOK, map[string]any{
		"message": "The Group post created with success",
		"GroupEvents":GetGroupEvents,
		"status":  "success",
	})
}
func (handler *Handlers) GroupEventParrticipate(res http.ResponseWriter, req *http.Request) {
}
