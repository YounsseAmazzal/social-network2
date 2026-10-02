package requests

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"social/internal/models"
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
