package handlers

import (
	"fmt"
	"net/http"
	"social/internal/utils"
)

func (*Handlers) SignUp(res http.ResponseWriter, req *http.Request) {
	fmt.Println("the user want to SignUp ")
	utils.ResponnseSender(res, http.StatusCreated, "User creaated succesfully ")
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
