package routers

import (
	"net/http"
	"social/internal/handlers"
	"social/internal/middleware"
	"social/internal/utils"
	"time"
)

func PooostsandComments(servermux *http.ServeMux, handler *handlers.Handlers, midware *middleware.Miiiddeware) {
	neewratelimiter := utils.NewRateLimiter(time.Minute)
	//alll auth routes
	servermux.Handle("/CreatePost", neewratelimiter.ApplyMiddleware("normalapi", (http.HandlerFunc(handler.CreatePosst))))
	servermux.Handle("/GetPost", neewratelimiter.ApplyMiddleware("normalapi", (http.HandlerFunc(handler.GETpost))))
	servermux.Handle("/AllPosts", neewratelimiter.ApplyMiddleware("normalapi", (http.HandlerFunc(handler.AllPosts))))
	servermux.Handle("/CreateComment", neewratelimiter.ApplyMiddleware("normalapi", (http.HandlerFunc(handler.CreateComment))))
	servermux.Handle("/GetComment", neewratelimiter.ApplyMiddleware("normalapi", (http.HandlerFunc(handler.CreateComment))))
	servermux.Handle("/AllComment", neewratelimiter.ApplyMiddleware("normalapi", (http.HandlerFunc(handler.AllComments))))
}
