package handlers

import (
	"encoding/json"
	"net/http"
	"social/internal/models"
	"social/internal/requests"
	"social/internal/utils"
)

func (handler *Handlers) SignUp(res http.ResponseWriter, req *http.Request) {
	//fmt.Println("the user want to SignUp ")
	if !requests.IsRequestValid(res, req, "/Signup", http.MethodPost) {
		return
	}
	// //fmt.Println("the resquet content",req.Header)
	// var userPayload models.UserPayload
	userpayload, err := requests.ValidSignUpPayload(req)
	if err != nil {
		//i will add a helper to send error
		// //fmt.Println("the payload mismached")
		utils.ErrorResponse(res, err.Error(), http.StatusBadRequest)
		return
	}
	// //fmt.Println("thaaat is the payload ConfirmPassword ", userpayload)
	if err := utils.ValidateRegistration(&userpayload); err != nil {
		// response errorrr
		// log.Fatal("probleeem on validation ",err.Error())
		utils.ErrorResponse(res, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	usersingup, err := handler.services.SignUpUser(userpayload)
	if err != nil {
		//fmt.Println("there is problem on user creations")
		// khaasni nspisefe error from database or what
		//i will add a hander error response
		utils.ErrorResponse(res, err.Error(), http.StatusBadRequest)
		return
	}
	errrrr := handler.services.Reepo.CreateUserSession(res, usersingup.ID)
	if errrrr != nil {
		//fmt.Println("i couldn't maaake session forr user")
		utils.ErrorResponse(res, err.Error(), http.StatusInternalServerError)
		return
		//add the handler error here
	}

	Userreturned := models.User{
		ID:        userpayload.ID,
		Nickname:  userpayload.Nickname,
		Username:  userpayload.Nickname,
		Firstname: userpayload.Firstname,
		Lastname:  userpayload.Lastname,
		Age:       userpayload.Age,
		Gender:    userpayload.Gender,
		Email:     userpayload.Email,
		AvatarURL: userpayload.AvatarURL,
		Avatar:    userpayload.AvatarURL,
		AboutMe:   userpayload.AboutMe,
		IsPrivate: userpayload.IsPrivate == 1,
	}
	//fmt.Println("that is the user payload", usersingup)
	utils.ResponnseSender(res, http.StatusCreated, map[string]any{
		"status":  "success",
		"message": "Welcome to the social network",
		"data": map[string]any{
			"user": Userreturned,
		},
	})
}
func (handler *Handlers) Login(res http.ResponseWriter, req *http.Request) {
	// //fmt.Println("the user want to Logout ")
	if !requests.IsRequestValid(res, req, "/Signin", http.MethodPost) {
		return
	}
	var userlogin models.UserLoign
	if err := json.NewDecoder(req.Body).Decode(&userlogin); err != nil {
		utils.ErrorResponse(res, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	//fmt.Println("the useeer crediantial", userlogin.IdentifierEN)
	err := utils.Logininputvalid(&userlogin)
	if err != nil {
		// log.Fatal("hnaaaaaaaaaaaaaaaaaaya")
		utils.ErrorResponse(res, err.Error(), http.StatusBadRequest)
		return
	}

	useeer, found := handler.services.Looign(userlogin.IdentifierEN)
	if !found || useeer == nil {
		utils.ErrorResponse(res, "Invaaalid IdentifierEN or password", http.StatusBadRequest)
		return
	}
	// //fmt.Println("that is the the passsword 1",useeer.Password)
	// //fmt.Println("that is the the passsword 2",userlogin.Password)
	if !utils.CompareHashedpassword(userlogin.Password, useeer.Password) {
		//fmt.Println("@@@@@@there is an problem on compeere @@@@@")
		utils.ResponnseSender(res, http.StatusCreated, "Invaaalid IdentifierEN or password")
		return
	}
	errrrr := handler.services.Reepo.CreateUserSession(res, useeer.ID)
	if errrrr != nil {
		//fmt.Println("i couldn't maaake session forr user")
		utils.ErrorResponse(res, err.Error(), http.StatusInternalServerError)
		return
		//add the handler error here
	}
	Userreturned := models.User{
		ID:         useeer.ID,
		Nickname:   useeer.Nickname,
		Username:   useeer.Nickname,
		Firstname:  useeer.Firstname,
		Lastname:   useeer.Lastname,
		Age:        useeer.Age,
		IsLoggedIn: true,
		Gender:     useeer.Gender,
		Email:      useeer.Email,
		AvatarURL:  useeer.AvatarURL,
		Avatar:     useeer.AvatarURL,
		AboutMe:    useeer.AboutMe,
		IsPrivate:  useeer.IsPrivate,
	}
	nothif:=models.Notification{
        Category: "Authontication",
		RecipientID: Userreturned.ID,
        Message: Userreturned.Nickname + " You logged successfully",
    }
	handler.services.Reepo.CreateNotification(&nothif)
	handler.wsmaanager.SendNotificationToUser(Userreturned.ID,nothif)
	utils.ResponnseSender(res, http.StatusCreated, map[string]any{
		"status":  "success",
		"message": "Welcome to the social network",
		"data": map[string]any{
			"user": Userreturned,
		},
	})
}
func (*Handlers) Home(res http.ResponseWriter, req *http.Request) {
	//fmt.Println("the user want to home ")
	utils.ResponnseSender(res, http.StatusCreated, "Yourrr welcome")
}
func (hander *Handlers) Logout(res http.ResponseWriter, req *http.Request) {
	// //fmt.Println("the user want to Logout ")
	if !requests.IsRequestValid(res, req, "/Logout", http.MethodDelete) {
		return
	}
	cookie, err := req.Cookie("solo_social")
	if err != nil {
		utils.ErrorResponse(res, "sesssion not valiiid ", http.StatusUnauthorized)
		return
	}
	_, err = hander.services.Reepo.Checksession(cookie.Value)
	if err != nil {
		utils.ErrorResponse(res, "Faileddd to delete the sesssions", http.StatusUnauthorized)
		// log.Printf("Loogout error",err.Error())
		return
	}
	err = hander.services.Reepo.DeleteSession(cookie.Value)
	if err != nil {
		utils.ErrorResponse(res, "Faileddd to delete the sesssions", http.StatusUnauthorized)
		return
	}
	http.SetCookie(res, &http.Cookie{
		Name:     "solo_social",
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
		Path:     "/",
		Secure:   false,
	})

	utils.ResponnseSender(res, http.StatusOK, map[string]string{
		"status":  "success",
		"message": "You Logout succesfully",
	})
}
func (handler *Handlers) ME(res http.ResponseWriter, req *http.Request) {
	//fmt.Println("the user want to HimSelf ")
	if !requests.IsRequestValid(res, req, "/Me", http.MethodGet) {
		return
	}
	coookie, err := req.Cookie("solo_social")
	if err != nil {
		utils.ErrorResponse(res, "sesssion not valiiid ", http.StatusUnauthorized)
		return
	}

	userID, err := handler.services.Reepo.Checksession(coookie.Value)
	if err != nil {
		utils.ErrorResponse(res, "NO Valid session exist", http.StatusUnauthorized)
		// log.Printf("Loogout error",err.Error())
		return
	}
	usserInfo, err := handler.services.Reepo.BringUserByid(userID)
	if err != nil {
		utils.ErrorResponse(res, "NO Valid session exist", http.StatusUnauthorized)
	}
	takeuser := models.User{
		ID:         usserInfo.ID,
		Nickname:   usserInfo.Nickname,
		Username:   usserInfo.Nickname,
		Firstname:  usserInfo.Firstname,
		Lastname:   usserInfo.Lastname,
		Age:        usserInfo.Age,
		IsLoggedIn: true,
		Gender:     usserInfo.Gender,
		Email:      usserInfo.Email,
		AvatarURL:  usserInfo.AvatarURL,
		Avatar:     usserInfo.AvatarURL,
		AboutMe:    usserInfo.AboutMe,
		IsPrivate:  usserInfo.IsPrivate,
	}

	utils.ResponnseSender(res, http.StatusOK, map[string]any{
		"message":  "GOT THe user Succefully",
		"Userinfo": takeuser,
	})
}
