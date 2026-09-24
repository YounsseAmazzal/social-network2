package routers

import (
	"net/http"
	"social/internal/handlers"
	"social/internal/middleware"
	"social/internal/utils"
	"time"
)

func ProfileAndFollowers(servermux *http.ServeMux, handler *handlers.Handlers, middwarre *middleware.Miiiddeware) {
	neewratelimiter := utils.NewRateLimiter(time.Minute)

	servermux.Handle("/UserProfile", neewratelimiter.ApplyMiddleware("normalapi", (http.HandlerFunc(handler.UserProfile))))
	servermux.Handle("/ProfileUppdate", neewratelimiter.ApplyMiddleware("normalapi", (http.HandlerFunc(handler.ProfileUppdate))))
	servermux.Handle("/followers/Accepting", neewratelimiter.ApplyMiddleware("normalapi", (http.HandlerFunc(handler.FollowersAccept))))
	servermux.Handle("/followers/Declineinng", neewratelimiter.ApplyMiddleware("normalapi", (http.HandlerFunc(handler.FollowersDeclineinng))))
	servermux.Handle("/followers/Pendiing", neewratelimiter.ApplyMiddleware("normalapi", (http.HandlerFunc(handler.FollowersPendiing))))
	servermux.Handle("/UserUnfollow", neewratelimiter.ApplyMiddleware("normalapi", (http.HandlerFunc(handler.UserUnfollow))))
	servermux.Handle("/Userfollowinng/", neewratelimiter.ApplyMiddleware("normalapi", (http.HandlerFunc(handler.Userfollowinng))))
	servermux.Handle("/Userfollowers/", neewratelimiter.ApplyMiddleware("normalapi", (http.HandlerFunc(handler.Userfollowers))))
}
