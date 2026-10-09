package handlers

import (
	"fmt"
	"net/http"
	"social/internal/requests"
	"social/internal/utils"
	"strconv"
)

func (handler *Handlers) Golbalenothification(res http.ResponseWriter, req *http.Request) {
	fmt.Println("The user want  All grouup messgess ")
	if !requests.IsRequestValid(res, req, "/golbalenothification", http.MethodGet) {
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
	Nothiflimit := requests.Nothiflimit(req)
	defaultlimit:=10
	notlimit, err := strconv.Atoi(Nothiflimit)
	if err!=nil{
		notlimit=defaultlimit
	}
	userNothification,err:=handler.services.Reepo.UserNothification(notlimit,user.ID)
	if err!=nil{
		utils.ErrorResponse(res, "Problem on getting nothifications", http.StatusInternalServerError)
		return
	}
	utils.ResponnseSender(res, http.StatusOK, map[string]any{
		"message":              "Taaake All posts",
		"status":               "success",
		"userNothification": userNothification,
	})

}
func (handler *Handlers) B(res http.ResponseWriter, req *http.Request) {
	
	utils.ResponnseSender(res, http.StatusOK, map[string]any{
		"message":              "Taaake All posts",
		"status":               "success",
	
	})

}
func (handler *Handlers) C(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Noooooooooooooooooooooooooooooooooooooooooooooooooooooooof")
	utils.ResponnseSender(res, http.StatusOK, map[string]any{
		"message":              "Taaake All posts",
		"status":               "success",
	
	})

}
func (handler *Handlers) D(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Noooooooooooooooooooooooooooooooooooooooooooooooooooooooof")
	
	utils.ResponnseSender(res, http.StatusOK, map[string]any{
		"message":              "Taaake All posts",
		"status":               "success",
	
	})

}