package routers

import (
	"net/http"
	"social/internal/handlers"
)

func WsChatAndNothification(servermux *http.ServeMux, wsm *handlers.WSManager) {
	servermux.HandleFunc("/WebSocket", wsm.WebSocket)
}
