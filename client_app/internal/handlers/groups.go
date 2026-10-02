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

	var createGroup models.GroupItem
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
		utils.ErrorResponse(res, err.Error(), http.StatusInternalServerError)
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
func (handler *Handlers) JoinedGroups(res http.ResponseWriter, req *http.Request) {
	if !requests.IsRequestValid(res, req, "/JoinedGroups", http.MethodGet) {
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
	joineddgroup, err := handler.services.Reepo.JooinedGroups(user.ID)
	if err != nil {
		utils.ErrorResponse(res, "Problem on bringing Joiined groups", http.StatusInternalServerError)
		return
	}
	utils.ResponnseSender(res, http.StatusOK, map[string]any{
		"status":       "success",
		"joineddgroup": joineddgroup,
	})

}
func (handler *Handlers) SeeeTheGroup(res http.ResponseWriter, req *http.Request) {
	if !requests.IsRequestValid(res, req, "/SeeeTheGroup", http.MethodGet) {
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
	GrouupId, err := requests.ExtraactGroupId(req)
	if err != nil {
		/// i'm working with the id from the parm so should notice that on front
		utils.ErrorResponse(res, "Gouuup Not avalible", http.StatusNotFound)
		return
	}
	Groupexist, err := handler.services.Reepo.GoupExist(GrouupId)
	if !Groupexist {
		utils.ErrorResponse(res, "Grouup Not Fouund", http.StatusNotFound)
		return
	}
	ismember, errr := handler.services.Reepo.CheckifUsserMember(user.ID, GrouupId)
	fmt.Println("The user is Member ?????", ismember)
	if errr != nil {
		utils.ErrorResponse(res, "Problem on checking membership", http.StatusNotFound)
		return
	}
	if !ismember {
		utils.ErrorResponse(res, "You  are not a member to see this group", http.StatusNotFound)
		return
	}
	groupDetail, err := handler.services.Reepo.GouupInfoo(user.ID, GrouupId)
	if err != nil {
		utils.ErrorResponse(res, "Problem on Bring Group details", http.StatusInternalServerError)
		return
	}
	utils.ResponnseSender(res, http.StatusOK, map[string]any{
		"status":       "success",
		"Groupdetails": groupDetail,
	})

}
func (handler *Handlers) GrouupMembers(res http.ResponseWriter, req *http.Request) {
	if !requests.IsRequestValid(res, req, "/GrouupMembers", http.MethodGet) {
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
	GrouupId, err := requests.ExtraactGroupId(req)
	if err != nil {
		/// i'm working with the id from the parm so should notice that on front
		utils.ErrorResponse(res, "Gouuup Not avalible", http.StatusNotFound)
		return
	}
	Groupexist, err := handler.services.Reepo.GoupExist(GrouupId)
	if !Groupexist {
		utils.ErrorResponse(res, "Grouup Not Fouund", http.StatusNotFound)
		return
	}
	ismember, errr := handler.services.Reepo.CheckifUsserMember(user.ID, GrouupId)
	fmt.Println("The user is Member ?????", ismember)
	if errr != nil {
		utils.ErrorResponse(res, "Problem on checking membership", http.StatusNotFound)
		return
	}
	if !ismember {
		utils.ErrorResponse(res, "You  are not a member to see this group", http.StatusNotFound)
		return
	}
	Groupsmembers, err := handler.services.Reepo.GrouupMembers(GrouupId)
	if err != nil {
		utils.ErrorResponse(res, "Problem on bringing Group Members", http.StatusInternalServerError)
		return
	}
	utils.ResponnseSender(res, http.StatusOK, map[string]any{
		"status":        "success",
		"Groupsmembers": Groupsmembers,
	})

}
func (handler *Handlers) CreatePostGroup(res http.ResponseWriter, req *http.Request) {
	if !requests.IsRequestValid(res, req, "/CreatePostGroup", http.MethodPost) {
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
	postPyyload, err := requests.ParseGroupPOst(req)
	err = requests.ValidGroupPost(&postPyyload)
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
	postPyyload.GroupPost.CreaterID = user.ID
	postPyyload.GrouupID = grouupid
	err = handler.services.Reepo.CreatePostGroup(&postPyyload)
	if err != nil {
		utils.ErrorResponse(res, "Problem on creating Group post try again ", http.StatusInternalServerError)
		// utils.ErrorResponse(res, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.ResponnseSender(res, http.StatusOK, map[string]string{
		"message": "The Group post created with success",
		"status":  "success",
	})
}
func (handler *Handlers) GroupPosts(res http.ResponseWriter, req *http.Request) {
	if !requests.IsRequestValid(res, req, "/GetGroupPosts", http.MethodGet) {
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

	// if err := json.NewDecoder(req.Body).Decode(&Grouppagiantion); err != nil {
	// 		utils.ErrorResponse(res, "Invalid Body Json", http.StatusBadRequest)
	// 	return
	// }

	Grouppagiantion := requests.GetTheLimitandOffset(req)
	fmt.Println("THe limit and the offset", Grouppagiantion.Limit)
	fmt.Println("THe limit and the offset", Grouppagiantion.Offset)
	grouupid, err := requests.ExtraactGroupId(req)
	if err != nil {
		utils.ErrorResponse(res, "Grouuup not Found", http.StatusUnauthorized)
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

	Groupexist, err := handler.services.Reepo.GoupExist(grouupid)
	if !Groupexist {
		utils.ErrorResponse(res, "Grouup Not Fouund", http.StatusNotFound)
		return
	}
	
	GroupPosts, err := handler.services.Reepo.GetAllGroupPosts(grouupid, Grouppagiantion.Limit, Grouppagiantion.Offset)
	if err != nil {
		utils.ErrorResponse(res, "Problem on Gettitng Group post try again ", http.StatusInternalServerError)
		// utils.ErrorResponse(res, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println("Check the offset ", len(GroupPosts))
	utils.ResponnseSender(res, http.StatusOK, map[string]any{
		"GroupPosts": GroupPosts,
		"status":     "success",
	})
}
func (handler *Handlers) GroupMembers(res http.ResponseWriter, req *http.Request) {
	fmt.Println("the user want to get the groups members ")
}
func (handler *Handlers) GroupCreateEvent(res http.ResponseWriter, req *http.Request) {
}
func (handler *Handlers) GroupEvents(res http.ResponseWriter, req *http.Request) {
}
func (handler *Handlers) GroupEventParrticipate(res http.ResponseWriter, req *http.Request) {
}
func (handler *Handlers) GroupInvite(res http.ResponseWriter, req *http.Request) {
}
