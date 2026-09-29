package handlers

import (
	"fmt"
	"net/http"
	"social/internal/requests"
	"social/internal/utils"
)

func (handler *Handlers) CreatePosst(res http.ResponseWriter, req *http.Request) {
	fmt.Println("The user want to create post")
	if !requests.IsRequestValid(res, req, "/CreatePost", http.MethodPost) {
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
	postPyyload,err:=requests.PostPyloadParse(req)
	err=requests.ValidPostPyload(&postPyyload)
	if err!=nil{
		fmt.Println("YAAAAAAAAAAAK",err)
		utils.ErrorResponse(res,err.Error(), http.StatusUnauthorized)
		return
	}
	postPyyload.CreaterID=user.ID
	err=handler.services.Reepo.CreaatePost(&postPyyload)
	if err!=nil{
		utils.ErrorResponse(res,"Problem on creating post try again ", http.StatusInternalServerError)
		return
	}
	utils.ResponnseSender(res,http.StatusOK,map[string]string{
		"message":"The post created with success",
		"status":"success",
	})
}

func (handler *Handlers) GETpost(res http.ResponseWriter, req *http.Request) {
	fmt.Println("The user want to Get GetPost")
	if !requests.IsRequestValid(res, req, "/GetPost", http.MethodGet) {
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
	Post_id,err:=requests.ExtractPosstid(req)
	if err!=nil{
		utils.ErrorResponse(res, "Poost Nooot Found ", http.StatusUnauthorized)
		return
	}
	kaynonGroup,err:=handler.services.Reepo.CheckPostOnGroupposts(Post_id)
	if kaynonGroup{
		/// posttt of grouuuup
	}
	post,errr:=handler.services.Reepo.GetCompletePostByID(Post_id,user.ID)
	if errr!=nil{
		utils.ErrorResponse(res,errr.Error(), http.StatusInternalServerError)
		return 
	}
		utils.ResponnseSender(res,http.StatusOK,map[string]any{
		"message":"Taaake All posts",
		"status":"success",
		"PostDetails":post,
	})
}
func (handler *Handlers) AllPosts(res http.ResponseWriter, req *http.Request) {
	fmt.Println("The user want  All posts")
	if !requests.IsRequestValid(res, req, "/AllPosts", http.MethodGet) {
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
	alllposts,err:=handler.services.Reepo.GetAllVisiblePosts(user.ID)
	if err!=nil{
			utils.ErrorResponse(res,err.Error(), http.StatusInternalServerError)
		return
	}
	utils.ResponnseSender(res,http.StatusOK,map[string]any{
		"message":"Taaake All posts",
		"status":"success",
		"FeedPosts":alllposts,
	})

}
func (handler *Handlers) CreateComment(res http.ResponseWriter, req *http.Request) {
	fmt.Println("The user want  to create comment")
	if !requests.IsRequestValid(res, req, "/CreateComment", http.MethodPost) {
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
	Post_id,err:=requests.ExtractPosstid(req)
	postOngroupGroup,err:=handler.services.Reepo.CheckPostOnGroupposts(Post_id)
	if err!=nil{
		utils.ErrorResponse(res,err.Error(), http.StatusInternalServerError)
		return 
	}
	if postOngroupGroup{
		/// posttt of grouuuup
		return
	}
	commentpayload,err:=requests.Parsecommmenpyload(req)
	err=requests.ValidCommentpy(&commentpayload)
	if err!=nil{
		utils.ErrorResponse(res,err.Error(), http.StatusUnauthorized)
		return
	}
	commentpayload.PostID=Post_id
	commentpayload.CreaterID=user.ID
	err=handler.services.Reepo.CreaatComment(&commentpayload)
	if err!=nil{
		utils.ErrorResponse(res,"Problem on creating Comments  try again ", http.StatusInternalServerError)
		return
	}
	utils.ResponnseSender(res,http.StatusOK,map[string]string{
		"message":"The Comment created with success",
		"status":"success",
	})
	// handler.services.CreatePostComment(res,req)
}
func CreateThepost(res *http.ResponseWriter,reeq *http.Request){

}
func (handler *Handlers) AllComments(res http.ResponseWriter, req *http.Request) {
	fmt.Println("The user want to Get post")
}
