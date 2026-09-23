package handlers

import (
	"fmt"
	"net/http"
	"social/internal/models"
	"social/internal/requests"
	"social/internal/utils"
)

func (handler *Handlers) SignUp(res http.ResponseWriter, req *http.Request) {
	fmt.Println("the user want to SignUp ")
	if !requests.IsRequestValid(res, req, "/Signup", http.MethodPost) {
		return
	}
	// var userPayload models.UserPayload
	userpayload, err := requests.ValidSignUpPayload(req)
	if err != nil {
		//i will add a helper to send error
		fmt.Println("the payload mismached")
		return
	}
	if err := utils.ValidateRegistration(&userpayload); err != nil {
		// response errorrr
		return
	}
	usersingup, err := handler.services.SignUpUser(userpayload)
	if err != nil {
		fmt.Println("there is problem on user creations")
		// khaasni nspisefe error from database or what 
		//i will add a hander error response
		return
	}
	errrrr:=handler.services.Reepo.CreateUserSession(res,userpayload.ID)
	if errrrr!=nil{
		fmt.Println("i couldn't maaake session forr user")
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
	fmt.Println("that is the user payload", usersingup)
    utils.ResponnseSender(res, http.StatusCreated, map[string]any{
        "status":  "success",
        "message": "Welcome to the social network",
        "data": map[string]any{
            "user": Userreturned,
        },
    })
}
func (*Handlers) Login(res http.ResponseWriter, req *http.Request) {
	fmt.Println("the user want to Logout ")
	utils.ResponnseSender(res, http.StatusCreated, "Your welcome again ")
}
func (*Handlers) Home(res http.ResponseWriter, req *http.Request) {
	fmt.Println("the user want to home ")
	utils.ResponnseSender(res, http.StatusCreated, "Yourrr welcome")
}
func (*Handlers) Logout(res http.ResponseWriter, req *http.Request) {
	fmt.Println("the user want to Logout ")
	utils.ResponnseSender(res, http.StatusCreated, "User loged succesfully ")
}
func (*Handlers) ME(res http.ResponseWriter, req *http.Request) {
	fmt.Println("the user want to Logout ")
	utils.ResponnseSender(res, http.StatusCreated, "User loged succesfully ")
}
