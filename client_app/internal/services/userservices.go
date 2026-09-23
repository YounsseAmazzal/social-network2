package services

import (
	"errors"
	"fmt"
	"social/internal/models"
	"social/internal/utils"
	"strings"
)

func (Service *Service) SignUpUser(userpayload models.UserInfo) (*models.User, error) {
	hasheedpassword, err := utils.HashPassword(userpayload.Password)
	if err != nil {
		return nil, errors.New("failed to secure user info")
	}
	userhowregisterd,err:=Service.Reepo.CreateUser(&userpayload)
	if err!=nil{
		if strings.Contains(err.Error(),"UNIQUE constraint failed"){
			return nil,errors.New("useername of email is already exissts ")
		}
	}
	fmt.Println("that is the password hasshed" ,hasheedpassword)
	return userhowregisterd, nil
}
