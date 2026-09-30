package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

func ResponnseSender(response http.ResponseWriter, staatuscode int, data interface{}) {
	response.Header().Set("Content-Type", "application/json")
	// dima
	response.WriteHeader(staatuscode)
	var errr error
	errr = json.NewEncoder(response).Encode(data)
	if errr != nil {
		log.Println("Errror fencoding Json", errr.Error())
	}
	//fmt.Println("THEEEE RESPONSE SENDED SUCCESFULLY")
}
func ErrorResponse(w http.ResponseWriter, msg string, staatusecode int) {
	errorResponse := map[string]string{"ERorrmessages": msg}
	ResponnseSender(w, staatusecode, errorResponse)
}
