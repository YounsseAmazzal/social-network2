package requests

import (
	"fmt"
	"net/http"
	"social/internal/models"
	"strings"
)

func ValidGroupcreation(grouppy *models.GroupItem) error {
	if len(grouppy.Description) > 1000 || len(grouppy.Description) == 0 {
		return fmt.Errorf("The discription must not be empty")
	}
	if len(grouppy.Title) > 1000 || len(grouppy.Title) == 0 {
		return fmt.Errorf("The Title must not be empty")
	}
	return nil

}
func ExtraactGroupId(request *http.Request) (string, error) {
	Groupid := request.URL.Query().Get("Grouup_id")
	if Groupid == "" {
		return "", fmt.Errorf("Grouup  Id not Foundd")
	}
	return Groupid, nil
}
func ExtraactInviationDecision(request *http.Request) (string, error) {
	decision := strings.ToLower(request.URL.Query().Get("type"))
	// fmt.Println("the grouuuu inviiiiite deciaotnsd fs=w",decision)
	if decision != "acceptinvitasion" && decision != "refuseinvitaion" {
		return "", fmt.Errorf("Grouup invete decision not Foundd")
	}
	return decision, nil
}
