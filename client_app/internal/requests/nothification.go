package requests

import (
	"net/http"
	"strings"
)

func Nothiflimit(request *http.Request) string {
	Nothiflimit := request.URL.Query().Get("Nothif_limit")
	if Nothiflimit == "" {
		Nothiflimit = "10"
	}
	return Nothiflimit
}
func ExtractReadParams(r *http.Request) (notificationID string, markAll bool) {
	id := r.URL.Query().Get("notification_id")
	allParam := r.URL.Query().Get("all")

	notificationID = strings.TrimSpace(id)

	markAll = (allParam == "true" || allParam == "1")
	
	return notificationID, markAll
}
