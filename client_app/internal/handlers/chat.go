package handlers

import (
	"fmt"
	"net/http"
	"social/internal/requests"
	"social/internal/utils"
)

func (handler *Handlers) Grouupmessages(res http.ResponseWriter, req *http.Request) {
	fmt.Println("The user want  All grouup messgess ")
	if !requests.IsRequestValid(res, req, "/Grouupmessages", http.MethodGet) {
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
	messagepg := requests.Messagepagination(req)
	fmt.Println("The lastmessage_time",messagepg.Last_time)
	fmt.Println("The liiiimit",messagepg.Limit)
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
	// fmt.Println("The user is Member ?????", ismember)
	if errr != nil {
		// fmt.Println("THE MEEEEMER CHECB err", errr)
		utils.ErrorResponse(res, "Problem on checking membership", http.StatusNotFound)
		return
	}
	if !ismember {
		utils.ErrorResponse(res, "You  are not a member to see this seeee the messages", http.StatusNotFound)
		return
	}
	Groupmessages, err := handler.services.Reepo.GetGroupMessagesCursor(grouupid, messagepg.Limit, messagepg.Last_time)
	if err != nil {
		utils.ErrorResponse(res, "Problem on getting Group messages", http.StatusInternalServerError)
		return
	}
	// fmt.Println("u recieve the messages succesfully ",Groupmessages)
	utils.ResponnseSender(res, http.StatusOK, map[string]any{
		"message":       "Taaake All posts",
		"status":        "success",
		"Groupmessages": Groupmessages,
	})

}
