package requests

import (
	"fmt"
	"net/http"
	"social/internal/models"
	"strconv"
)

func Messagepagination(request *http.Request) models.Messagespg {
	liimit := request.URL.Query().Get("limit")
	last_time := request.URL.Query().Get("last_time")
	limit, err := strconv.Atoi(liimit)
	if err != nil {
		limit = 10
	}
	Messagepagination := models.Messagespg{Limit: limit, Last_time: last_time}
	Messagepagination.SetMessageDef()
	return Messagepagination
}
func ExtractUsertoMessage(request *http.Request) (string, error) {
	InvetedUser := request.URL.Query().Get("Userchat_id")
	if InvetedUser == "" {
		return "", fmt.Errorf("Userchat_id Id not Foundd")
	}
	return InvetedUser, nil
}
