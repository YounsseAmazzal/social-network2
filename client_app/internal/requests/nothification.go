package requests

import (
	"net/http"
)

func Nothiflimit(request *http.Request) (string) {
	Nothiflimit := request.URL.Query().Get("Nothif_limit")
	if Nothiflimit == "" {
		Nothiflimit = "10"
	}
	return Nothiflimit
}
