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
	servermux.Handle("/UsersForChat", neewratelimiter.ApplyMiddleware("normalapi", (http.HandlerFunc(handler.UsersForChat))))
	servermux.Handle("/Usersmessages", neewratelimiter.ApplyMiddleware("normalapi", (http.HandlerFunc(handler.Usersmessages))))
	servermux.Handle("/golbalenothification", neewratelimiter.ApplyMiddleware("normalapi", (http.HandlerFunc(handler.Golbalenothification))))
	servermux.Handle("/notifications/chat", neewratelimiter.ApplyMiddleware("normalapi", (http.HandlerFunc(handler.C))))
	servermux.Handle("/notifications/groups", neewratelimiter.ApplyMiddleware("normalapi", (http.HandlerFunc(handler.B))))

}
