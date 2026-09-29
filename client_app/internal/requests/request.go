package requests

import (
	"log"
	"net/http"
	"social/internal/utils"
	"strings"
)

// insure that the method and the path is good
func IsRequestValid(w http.ResponseWriter, r *http.Request, targetPath, allowedMethod string) bool {
	if strings.Contains(targetPath, "*") {
		pathSegments := strings.Split(r.URL.Path, "/")
		if len(pathSegments) > 0 {
			lastSegment := pathSegments[len(pathSegments)-1]
			targetPath = targetPath[:len(targetPath)-1] + lastSegment
		}
	}
	if r.URL.Path != targetPath {
		// HandleError(w, http.StatusNotFound, "Page not found")
		utils.ErrorResponse(w, "Page not found", http.StatusNotFound)
		log.Printf("404 Path mismatch")
		return false
	}
	if r.Method != allowedMethod {
		w.WriteHeader(http.StatusMethodNotAllowed)
		log.Printf("405 Method Not Allowed")
		utils.ErrorResponse(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return false
	}
	return true
}
