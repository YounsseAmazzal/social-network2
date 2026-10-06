package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"social/internal/models"
	"social/internal/requests"
	"social/internal/utils"
	"strings"
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
func (handler *Handlers) CreateGPComment(res http.ResponseWriter, req *http.Request) {
	if !requests.IsRequestValid(res, req, "/CreateGPComment", http.MethodPost) {
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
	Commentpy, err := requests.Parsecommmenpyload(req)
	err = requests.ValidCommentpy(&Commentpy)
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
		utils.ErrorResponse(res, "You  are not a member to see this group to comment", http.StatusNotFound)
		return
	}
	Post_id, err := requests.ExtractPosstid(req)
	if err != nil {
		utils.ErrorResponse(res, "Poost Nooot Found ", http.StatusUnauthorized)
		return
	}
	kaynonGroup, err := handler.services.Reepo.CheckPostOnGroupposts(Post_id)
	if kaynonGroup {
		// fmt.Println("COmmmment on group")
		Commentpy.CreaterID = user.ID
		Commentpy.PostID = Post_id
		err = handler.services.Reepo.CreateGroupComment(&Commentpy)
		if err != nil {
			utils.ErrorResponse(res, "Problem on creation Group Post Comment try again  ", http.StatusInternalServerError)
			return
		}
		utils.ResponnseSender(res, http.StatusOK, map[string]string{
			"message": "The Group Comment Post Created success",
			"status":  "success",
		})
		return
	}
	utils.ErrorResponse(res, "somthing went wrong ....", http.StatusInternalServerError)
}

func (handler *Handlers) GroupInvite(res http.ResponseWriter, req *http.Request) {
}
func (handler *Handlers) RequesttoJoinGroup(res http.ResponseWriter, req *http.Request) {
	if !requests.IsRequestValid(res, req, "/RequesttoJoinGroup", http.MethodPost) {
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
	// fmt.Println("The user is Member ?????", ismember)
	if errr != nil {
		// fmt.Println("THE MEEEEMER CHECB err", errr)
		utils.ErrorResponse(res, "Problem on checking membership try again", http.StatusInternalServerError)
		return
	}
	if ismember {
		utils.ErrorResponse(res, "You aree a member on this group", http.StatusForbidden)
		return
	}
	// log.Fatal("he is a member", ismember)
	err = handler.services.Reepo.GroupJoinRequests(grouupid, user.ID)
	if err != nil {
		utils.ErrorResponse(res, "Problem on requesting to join the group", http.StatusInternalServerError)
		return
	}
	utils.ResponnseSender(res, http.StatusOK, map[string]string{
		"message": "You requet to join the group successfully",
		"status":  "success",
	})
}
func (handler *Handlers) JoiningGroupRequests(res http.ResponseWriter, req *http.Request) {
	if !requests.IsRequestValid(res, req, "/JoiningGroupRequests", http.MethodGet) {
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
	groupadmin, err := handler.services.Reepo.GetGroupAdmin(grouupid)
	if err != nil {
		utils.ErrorResponse(res, "Grouup not found  ", http.StatusNotFound)
		return
	}
	if groupadmin != user.ID {
		utils.ErrorResponse(res, "You are not the admin of the grouup ", http.StatusUnauthorized)
		return
	}
	///
	pendingList, err := handler.services.Reepo.GetJoinRequesters(grouupid)
	if err != nil {
		utils.ErrorResponse(res, "Faailed to bring the requesters list ", http.StatusUnauthorized)
		return
	}
	utils.ResponnseSender(res, http.StatusOK, map[string]any{
		"status":        "success",
		"JoinGrouplist": pendingList,
	})
}
func (handler *Handlers) RespondtoJoinRequest(res http.ResponseWriter, req *http.Request) {
	if !requests.IsRequestValid(res, req, "/RespondtoJoinRequest", http.MethodPost) {
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
	decision := strings.ToLower(req.URL.Query().Get("type"))
	fmt.Println("the decision", decision)
	if decision != "joinaccept" && decision != "joinreject" {
		utils.ErrorResponse(res, "Invalid type. Must be 'Joinaccept' or 'Joinreject'", http.StatusBadRequest)
		return
	}
	applicationId := strings.ToLower(req.URL.Query().Get("useeer_id"))
	if applicationId == "" {
		utils.ErrorResponse(res, "Invalid User application id. ", http.StatusBadRequest)
		return
	}
	userexxist, err := handler.services.Reepo.UserExistsByID(applicationId)
	if err != nil {
		utils.ErrorResponse(res, "Database Errror", http.StatusInternalServerError)
		return
	}
	if !userexxist {
		utils.ErrorResponse(res, "User that requested no exist", http.StatusNotFound)
		return
	}

	groupadmin, err := handler.services.Reepo.GetGroupAdmin(grouupid)
	if err != nil {
		utils.ErrorResponse(res, "Grouup not found  ", http.StatusNotFound)
		return
	}
	if groupadmin != user.ID {
		utils.ErrorResponse(res, "only the admin can take the discion of acceptation ", http.StatusUnauthorized)
		return
	}
	fmt.Println("The group id ", grouupid)
	fmt.Println("THe user request id ", applicationId)
	isPending, err := handler.services.Reepo.IsJoinRequestPending(grouupid, applicationId)
	if err != nil {
		utils.ErrorResponse(res, "Database pending Errror", http.StatusInternalServerError)
		return
	}
	if !isPending {
		utils.ErrorResponse(res, "No pending requests avaailiable from this user", http.StatusNotFound)
		return
	}
	err = handler.services.Reepo.JoinRequestProcces(grouupid, applicationId, decision)
	if err != nil {
		utils.ErrorResponse(res, "Database accept and reject Errror", http.StatusInternalServerError)
		return
	}
	// i need nothification now here to optimaze the work
	utils.ResponnseSender(res, http.StatusOK, map[string]any{
		"status":  "success",
		"message": "You take the join decions succefully  ",
	})
}
func (handler *Handlers) InviteUserstogroup(res http.ResponseWriter, req *http.Request) {
	if !requests.IsRequestValid(res, req, "/InviteUserstogroup", http.MethodGet) {
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
	// fmt.Println("The user is Member ?????", ismember)
	if errr != nil {
		// fmt.Println("THE MEEEEMER CHECB err", errr)
		utils.ErrorResponse(res, "Problem on checking membership try again", http.StatusInternalServerError)
		return
	}
	if !ismember {
		utils.ErrorResponse(res, "You aree not member on this group", http.StatusForbidden)
		return
	}
	// log.Fatal("he is a member", ismember)
	err = handler.services.Reepo.GroupJoinRequests(grouupid, user.ID)
	if err != nil {
		utils.ErrorResponse(res, "Problem on requesting to join the group", http.StatusInternalServerError)
		return
	}
	InvidetUserList, err := handler.services.Reepo.GetFollowersToinvite(grouupid, user.ID)
	if err != nil {
		utils.ErrorResponse(res, "Problem on Getting followers to invite them", http.StatusInternalServerError)
		return
	}
	utils.ResponnseSender(res, http.StatusOK, map[string]any{
		"status":          "success",
		"message":         "You got the user invite list successfuly",
		"InvidetUserList": InvidetUserList,
	})
}
func (handler *Handlers) InviteGroup(res http.ResponseWriter, req *http.Request) {
	if !requests.IsRequestValid(res, req, "/InviteGroup", http.MethodPost) {
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
	// fmt.Println("The user is Member ?????", ismember)
	if errr != nil {
		// fmt.Println("THE MEEEEMER CHECB err", errr)
		utils.ErrorResponse(res, "Problem on checking membership try again", http.StatusInternalServerError)
		return
	}
	if !ismember {
		utils.ErrorResponse(res, "You aree not member on this group", http.StatusForbidden)
		return
	}
	inveteeduser, err := requests.Invetuserextracct(req)
	if err != nil {
		utils.ErrorResponse(res, err.Error(), http.StatusNotFound)
		return
	}
	if inveteeduser == user.ID {
		utils.ErrorResponse(res, "You can Not invete you Sefl", http.StatusNotFound)
		return
	}
	inveteeismember, errr := handler.services.Reepo.CheckifUsserMember(user.ID, grouupid)
	// fmt.Println("The user is Member ?????", ismember)
	if errr != nil {
		// fmt.Println("THE MEEEEMER CHECB err", errr)
		utils.ErrorResponse(res, "Problem on checking membership try again", http.StatusInternalServerError)
		return
	}
	if !inveteeismember {
		utils.ErrorResponse(res, "The user is already memeber", http.StatusForbidden)
		return
	}
	err=handler.services.Reepo.SendInvetation(grouupid,user.ID,inveteeduser)
		if err!=nil {
		utils.ErrorResponse(res, "Error Send the Invetion of group", http.StatusForbidden)
		return
		}
	utils.ResponnseSender(res, http.StatusOK, map[string]any{
		"status":  "success",
		"message": "You got the user invite list successfuly",
	})
}
