package requests

import (
	"fmt"
	"net/http"
)

func Invetuserextracct(request *http.Request) (string, error) {
	InvetedUser := request.URL.Query().Get("InvetUser_id")
	if InvetedUser == "" {
		return "", fmt.Errorf("InvetedUser Id not Foundd")
	}
	return InvetedUser, nil
}

