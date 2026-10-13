package handlers

import (
	"fmt"
	"log"
	"net/http"
	"social/internal/requests"
	"social/internal/utils"
)

func (handler *Handlers) CreatePosst(res http.ResponseWriter, req *http.Request) {
	//fmt.Println("The user want to create post")
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
	postPyyload, err := requests.PostPyloadParse(req)
	// log.Fatal("the pic",postPyyload.Image)
	err = requests.ValidPostPyload(&postPyyload)
	if err != nil {
		//fmt.Println("YAAAAAAAAAAAK", err)
		utils.ErrorResponse(res, err.Error(), http.StatusUnauthorized)
		return
	}
	postPyyload.CreaterID = user.ID
	err = handler.services.Reepo.CreaatePost(&postPyyload)
	if err != nil {
		utils.ErrorResponse(res, "Problem on creating post try again ", http.StatusInternalServerError)
		return
	}
	utils.ResponnseSender(res, http.StatusOK, map[string]string{
		"message": "The post created with success",
		"status":  "success",
	})
}

func (handler *Handlers) GETpost(res http.ResponseWriter, req *http.Request) {
	//fmt.Println("The user want to Get GetPost")
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
	Post_id, err := requests.ExtractPosstid(req)
	if err != nil {
		utils.ErrorResponse(res, "Poost Nooot Found ", http.StatusUnauthorized)
		return
	}
	kaynonGroup, err := handler.services.Reepo.CheckPostOnGroupposts(Post_id)
	if kaynonGroup {
		// fmt.Println("Post EXISTS on grouuuuup post ,")

		postgroup, err := handler.services.Reepo.PostGrouupDetailed(user.ID, Post_id)
		if err != nil {
			fmt.Println("EDDDDDDDDDD")
			utils.ErrorResponse(res, "Probem to serve post or you are not member", http.StatusUnauthorized)
			return
		}
			Grouppagiantion := requests.GetTheLimitandOffset(req)
			// fmt.Println("THe limit and the offset$$$$$$$$$$", Grouppagiantion.Limit)
			// fmt.Println("THe limit and the offset$$$$$$$$$$", Grouppagiantion.Offset)
		GrouppostComments, err := handler.services.Reepo.PostGroupComments(user.ID, Post_id,Grouppagiantion.Limit,Grouppagiantion.Offset)
		if err != nil {
			// fmt.Println("HHEHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHHH")
			utils.ErrorResponse(res, "Probem to serve post or you are not member", http.StatusUnauthorized)
			return
		}
		// fmt.Println("@@@@@@@@@@@@@@@@@@@@@@POstttt Commments@@@@@@@@@@@@@@@")
		// fmt.Println("@@@", GrouppostComments)
		// fmt.Println("@@@@@@@@@@@@@@@@@@@@@@POstttt Commments@@@@@@@@@@@@@@@")
		// fmt.Println("@@@@@@@@@@@@@@@@@@@@@@POstttt Commments@@@@@@@@@@@@@@@")

		postgroup.PostComments = GrouppostComments
		utils.ResponnseSender(res, http.StatusOK, map[string]any{
			"message":           "Get group post detailed",
			"status":            "success",
			"PostGroupDetailed": postgroup,
		})
		/// posttt of grouuuup
		return
	}
	post, errr := handler.services.Reepo.GetCompletePostByID(Post_id, user.ID)
	if errr != nil {
		utils.ErrorResponse(res, errr.Error(), http.StatusInternalServerError)
		return
	}
	Postcommmentp := requests.PostPagination(req)
			// fmt.Println("THe limit and the offset$$$$$$$$$$", Postcommmentp.Limit)
			// fmt.Println("THe limit and the offset$$$$$$$$$$", Postcommmentp.Offset)
	postcomments, err := handler.services.Reepo.GetPostCommensts(user.ID, Post_id,Postcommmentp.Limit,Postcommmentp.Offset)
	if err != nil {
		utils.ErrorResponse(res, "WE HAVE A trouble try again", http.StatusInternalServerError)
		return
	}
	post.PostComments = postcomments
	//fmt.Println("PPPPPOOOOOST COMMMENTS",postcomments)
	utils.ResponnseSender(res, http.StatusOK, map[string]any{
		"message":     "Taaake All posts",
		"status":      "success",
		"PostDetails": post,
	})
}
func (handler *Handlers) LikePost(res http.ResponseWriter, req *http.Request) {
	if !requests.IsRequestValid(res, req, "/Postlike", http.MethodPost) {
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
	Post_id, err := requests.ExtractPosstid(req)
	// //fmt.Println("::::::::::::::::::::::::::::::::::::::")
	// //fmt.Println("::::::::::::::::::::::::::::::::::::::")
	// //fmt.Println("::::::::::::::::::::::::::::::::::::::")
	// //fmt.Println(":::::::::::::post id:::::::::::::")
	// //fmt.Println("::::::::::::::::::::::::::::::::::::::",Post_id)
	// //fmt.Println("::::::::::::::::::::::::::::::::::::::")
	if err != nil {
		utils.ErrorResponse(res, "Poost Nooot Found ", http.StatusNotFound)
		return
	}

	// userhasrate,rateexit,err:=handler.services.Reepo.GetPostrating(Post_id,user.ID)
	// if err!=nil{
		// 	utils.ErrorResponse(res, "Problem on raating sytem ", http.StatusUnauthorized)
		// 	return
		// }
	kaynonGroup, err := handler.services.Reepo.CheckPostOnGroupposts(Post_id)
	if kaynonGroup {
			// fmt.Println("The post exxist on group posts ")
				err = handler.services.Reepo.GroupPostLike(1, user.ID, Post_id)
				if err != nil {
					utils.ErrorResponse(res, err.Error(), http.StatusInternalServerError)
					return
				}
					utils.ResponnseSender(res, http.StatusOK, map[string]any{
							"message": "group post liked on success",
							"status":  "success",
					})
			/// posttt of grouuuup
			return
	}
	postexist, err := handler.services.Reepo.CheeeckPostidvalid(Post_id)
	if !postexist {
		utils.ErrorResponse(res, "POST NOT FOUND", http.StatusNotFound)
		return
	}
	err = handler.services.Reepo.PostLike(1, user.ID, Post_id)
	if err != nil {
		utils.ErrorResponse(res, err.Error(), http.StatusInternalServerError)
		return
	}
	// if userhasrate{
	// 	////send notification
	// }
	utils.ResponnseSender(res, http.StatusOK, map[string]any{
		"message": "post liked on success",
		"status":  "success",
	})

}
func (handler *Handlers) DislikePost(res http.ResponseWriter, req *http.Request) {
	if !requests.IsRequestValid(res, req, "/Postdislike", http.MethodPost) {
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
	Post_id, err := requests.ExtractPosstid(req)
	if err != nil {
		utils.ErrorResponse(res, "Poost Nooot Found ", http.StatusUnauthorized)
		return
	}
	kaynonGroup, err := handler.services.Reepo.CheckPostOnGroupposts(Post_id)
	if kaynonGroup {
			// fmt.Println("The post exxist on group posts ")
				err = handler.services.Reepo.GroupPostLike(0, user.ID, Post_id)
				if err != nil {
					utils.ErrorResponse(res, err.Error(), http.StatusInternalServerError)
					return
				}
					utils.ResponnseSender(res, http.StatusOK, map[string]any{
							"message": "group post disliked on success",
							"status":  "success",
					})
			/// posttt of grouuuup
			return
	}
	postexist, err := handler.services.Reepo.CheeeckPostidvalid(Post_id)
	if !postexist {
		utils.ErrorResponse(res, "POST NOT FOUND ON DISLIKE ", http.StatusUnauthorized)
		return
	}
	// userhasrate,rateexit,err:=handler.services.Reepo.GetPostrating(Post_id,user.ID)
	// if err!=nil{
	// 	utils.ErrorResponse(res, "Problem on raating sytem ", http.StatusUnauthorized)
	// 	return
	// }
	err = handler.services.Reepo.PostLike(0, user.ID, Post_id)
	if err != nil {
		utils.ErrorResponse(res, err.Error(), http.StatusInternalServerError)
		return
	}
	// if userhasrate{
	// 	////send notification
	// }
	utils.ResponnseSender(res, http.StatusOK, map[string]any{
		"message": "post disliked on success",
		"status":  "success",
	})

}
func (handler *Handlers) AllPosts(res http.ResponseWriter, req *http.Request) {
	//fmt.Println("The user want  All posts")
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
	PostPagination := requests.PostPagination(req)
	// fmt.Println("THe limit and the offset", PostPagination.Limit)
	// fmt.Println("THe limit and the offset", PostPagination.Offset)
	alllposts, err := handler.services.Reepo.GetAllVisiblePosts(user.ID,PostPagination.Limit,PostPagination.Offset)
	if err != nil {
		utils.ErrorResponse(res, err.Error(), http.StatusInternalServerError)
		return
	}
	utils.ResponnseSender(res, http.StatusOK, map[string]any{
		"message":   "Taaake All posts",
		"status":    "success",
		"FeedPosts": alllposts,
	})

}
func (handler *Handlers) CreateComment(res http.ResponseWriter, req *http.Request) {
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
	Post_id, err := requests.ExtractPosstid(req)
	// //fmt.Println("THeeeeeeeee post id of the comment",Post_id)
	if err != nil {
		utils.ErrorResponse(res, err.Error(), http.StatusInternalServerError)
		return
	}
	
	commentpayload, err := requests.Parsecommmenpyload(req)
	err = requests.ValidCommentpy(&commentpayload)
	if err != nil {
		utils.ErrorResponse(res, err.Error(), http.StatusUnauthorized)
		return
	}
	postOngroupGroup, err := handler.services.Reepo.CheckPostOnGroupposts(Post_id)
	
	commentpayload.PostID = Post_id
	commentpayload.CreaterID = user.ID
	// log.Fatal("YAAAAAAAAAAAAAAAAAAAAAAAAK ")
	if postOngroupGroup {
		err = handler.services.Reepo.CreateGroupComment(&commentpayload)
		if err!=nil{
			log.Fatal("YAAAAAAAAAAAAAAAAAAAAAK ")
				utils.ErrorResponse(res, err.Error(), http.StatusUnauthorized)
			return
		}
		// log.Fatal("commment created")
		utils.ResponnseSender(res, http.StatusOK, map[string]string{
		"message": "The Comment Post Group created with success",
		"status":  "success",
		})
		/// posttt of grouuuup
		return
	}
	postexist, err := handler.services.Reepo.CheeeckPostidvalid(Post_id)
	if !postexist {
		utils.ErrorResponse(res, "POST NOT FOUND", http.StatusUnauthorized)
		return
	}
	err = handler.services.Reepo.CreaatComment(&commentpayload)
	if err != nil {
		utils.ErrorResponse(res, "Problem on creating Comments  try again ", http.StatusInternalServerError)
		return
	}
	utils.ResponnseSender(res, http.StatusOK, map[string]string{
		"message": "The Comment created with success",
		"status":  "success",
	})
	// handler.services.CreatePostComment(res,req)
}


// func (handler *Handlers) GetComments(res http.ResponseWriter, req *http.Request) {
// 	// //fmt.Println("The user want  to create comment")
// 	if !requests.IsRequestValid(res, req, "/AllPostComment", http.MethodGet) {
// 		return
// 	}
// 	cookie, err := req.Cookie("solo_social")
// 	if err != nil {
// 		utils.ErrorResponse(res, "sesssion not valiiid ", http.StatusUnauthorized)
// 		return
// 	}
// 	_, err = handler.services.Reepo.Checksession(cookie.Value)
// 	if err != nil {
// 		utils.ErrorResponse(res, "Your Session Not valid", http.StatusUnauthorized)
// 		// log.Printf("Your Session Not valid", err.Error())
// 		return
// 	}
// 	// user, err := handler.services.Reepo.GetUserByToken(cookie.Value)
// 	// if err != nil {
// 	// 	utils.ErrorResponse(res, "Your Session Not validd", http.StatusUnauthorized)
// 	// 	return
// 	// }
// 	// handler.services.Reepo.
// 		utils.ResponnseSender(res, http.StatusOK, map[string]string{
// 		"message": "The Comment created with success",
// 		"status":  "success",
// 	})
// 	// handler.services.CreatePostComment(res,req)
// }
// func CreateThepost(res *http.ResponseWriter, reeq *http.Request) {

// }
// func (handler *Handlers) GetPostComment(res http.ResponseWriter, req *http.Request) {
// 	//fmt.Println("The user want to Get GetPostComment")
// }
