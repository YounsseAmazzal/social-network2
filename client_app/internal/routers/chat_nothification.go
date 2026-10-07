package routers

import (
	"net/http"
	"social/internal/handlers"
	"social/internal/middleware"
	"social/internal/utils"
	"time"
)

func WsChatAndNothification(servermux *http.ServeMux, wsm *handlers.WSManager, middleware *middleware.Miiiddeware,handler *handlers.Handlers) {
	servermux.HandleFunc("/WebSocket", wsm.WebSocket)
	// groouup messages
	neewratelimiter := utils.NewRateLimiter(time.Minute)
	servermux.Handle("/Grouupmessages", neewratelimiter.ApplyMiddleware("normalapi", (http.HandlerFunc(handler.Grouupmessages))))

}
