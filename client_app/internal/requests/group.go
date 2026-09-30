package requests

import (
	"fmt"
	"social/internal/models"
)
func ValidGroupcreation(grouppy *models.Group)error{	
		if len(grouppy.Description) > 1000 || len(grouppy.Description)==0{
			return  fmt.Errorf("The discription must not be empty")
		}
		if len(grouppy.Title) > 1000 || len(grouppy.Title)==0{
			return  fmt.Errorf("The Title must not be empty")
		}
		return  nil
	
}