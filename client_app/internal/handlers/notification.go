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
func (handler *Handlers) ChatNothification(res http.ResponseWriter, req *http.Request) {
	if !requests.IsRequestValid(res, req, "/ChatNothification", http.MethodGet) {
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
	counters,err:=handler.services.Reepo.CountermesaageNotRead(user.ID)
	if err != nil {
        utils.ErrorResponse(res,"Database error",500)
        return
    }
	unreadMap := make(map[string]int)
    totalUnread := 0

    for _, c := range counters {
		// for everrry userid [useeif : 44]
		unreadMap[c.SenderID] = c.Count
        totalUnread += c.Count
    }

	utils.ResponnseSender(res, http.StatusOK, map[string]any{
		"message":              "Taaake All Nothiiif",
		"status":               "success",
		"senders":counters,
		"NotRead_Total":totalUnread,
	})

}
func (handler *Handlers) GroupChatNothification(res http.ResponseWriter, req *http.Request) {
	if !requests.IsRequestValid(res, req, "/GroupChatNothification", http.MethodGet) {
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
	// i think i don't need to check the membersheeeep 

	// GrouupId, err := requests.ExtraactGroupId(req)
	// if err != nil {
	// 	/// i'm working with the id from the parm so should notice that on front
	// 	utils.ErrorResponse(res, "Gouuup Not avalible", http.StatusNotFound)
	// 	return
	// }
	// Groupexist, err := handler.services.Reepo.GoupExist(GrouupId)
	// if !Groupexist {
	// 	utils.ErrorResponse(res, "Grouup Not Fouund", http.StatusNotFound)
	// 	return
	// }
	// ismember, errr := handler.services.Reepo.CheckifUsserMember(user.ID, GrouupId)
	// fmt.Println("The user is Member ?????", ismember)
	// if errr != nil {
	// 	utils.ErrorResponse(res, "Problem on checking membership", http.StatusNotFound)
	// 	return
	// }
	// if !ismember {
	// 	utils.ErrorResponse(res, "You  are not a member to see this group", http.StatusNotFound)
	// 	return
	// }
	// groupDetail, err := handler.services.Reepo.GouupInfoo(user.ID, GrouupId)
	// if err != nil {
	// 	utils.ErrorResponse(res, "Problem on Bring Group details", http.StatusInternalServerError)
	// 	return
	// }

	counters,err:=handler.services.Reepo.CounterGroupmesaageNotRead(user.ID)
	if err != nil {
        utils.ErrorResponse(res,"Database error",500)
        return
    }
    groupCounts := make(map[string]int)
    totalGroupUnread := 0

    for _, c := range counters {
        groupCounts[c.GroupID] = c.Count
        totalGroupUnread += c.Count
    }


	utils.ResponnseSender(res, http.StatusOK, map[string]any{
		"message":              "Taaake All Nothiiif",
		"status":               "success",
		"senders":counters,
		"NotRead_Total":groupCounts,
	})

}
func (handler *Handlers) ReadGroupNothification(res http.ResponseWriter, req *http.Request) {
	if !requests.IsRequestValid(res, req, "/ReadGroupNothification", http.MethodPost) {
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
	grouupid,err:=requests.ExtraactGroupId(req)
	if err!=nil ||grouupid==""{
		utils.ErrorResponse(res, "Grouup Not Found", http.StatusUnauthorized)
		return
	}
	handler.services.Reepo.ReadNgroupnothif(user.ID,grouupid)
	utils.ResponnseSender(res, http.StatusOK, map[string]any{
		"message":              "Taaake All Nothiiif",
		"status":               "success",
	})

}
func (handler *Handlers) ReadNothification(res http.ResponseWriter, req *http.Request) {
	if !requests.IsRequestValid(res, req, "/ReadNothification", http.MethodPost) {
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
	talker,err:=requests.ExtractUsertoMessage(req)
	if err != nil {
		utils.ErrorResponse(res, "Missing talkeer ", http.StatusNotFound)
		return
	}
	err=handler.services.Reepo.Readchatnothif(user.ID,talker)
	if err!=nil{
			utils.ErrorResponse(res, "ERror on marking chat messages nothification as read", http.StatusNotFound)
		return
	}
	// counters,err:=handler.services.Reepo.CounterGroupmesaageNotRead(user.ID)
	// if err != nil {
    //     utils.ErrorResponse(res,"Database error",500)
    //     return
    // }
    // groupCounts := make(map[string]int)
    // totalGroupUnread := 0

    // for _, c := range counters {
    //     groupCounts[c.GroupID] = c.Count
    //     totalGroupUnread += c.Count
    // }


	utils.ResponnseSender(res, http.StatusOK, map[string]any{
		"message":              " chat nothification messages marked as read",
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