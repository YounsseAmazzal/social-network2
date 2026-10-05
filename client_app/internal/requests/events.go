package requests

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"social/internal/models"
	"strconv"
	"strings"
	"time"
)

func Eventplyload(req *http.Request) (models.Eventplyload, error) {
	var Eventplyload models.Eventplyload
	if err := json.NewDecoder(req.Body).Decode(&Eventplyload); err != nil {
		return Eventplyload, fmt.Errorf("Invalid json form")
	}
	return Eventplyload, nil
}

func ExtractEventID(request *http.Request) (string, error) {
	Groupid := request.URL.Query().Get("Event_id")
	if Groupid == "" {
		return "", fmt.Errorf("Event_id  Id not Foundd")
	}
	return Groupid, nil
}

func Vaalidevent(event *models.Eventplyload) error {
	if strings.TrimSpace(event.Title) == "" || len(event.Title) > 100 {
		return errors.New("title is required and must be under 100 characters")
	}
	if event.EventDate.Before(time.Now()) {
		return errors.New("event date cannot be in the past")
	}
	if event.EventDate.After(time.Now().AddDate(5, 0, 0)) {
		return errors.New("event cannot be scheduled more than 5 years in advance")
	}
	return nil
}
func ExtractEventVote(r *http.Request) string {
	vote := r.URL.Query().Get("vote")

	vote = strings.ToLower(strings.TrimSpace(vote))

	if vote != "going" && vote != "not going" && vote != "remove" {
		return ""
	}

	return vote
}

func EventPagination(request *http.Request) models.EventPagination {
	liimit := request.URL.Query().Get("limit")
	offsetstr := request.URL.Query().Get("offset")
	limit, err := strconv.Atoi(liimit)
	if err != nil {
		limit = 10
	}
	offset, err := strconv.Atoi(offsetstr)
	if err != nil {
		offset = 0
	}
	EventPagination := models.EventPagination{Limit: limit, Offset: offset}
	EventPagination.SetEventDefaults()
	return EventPagination
}
