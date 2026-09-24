package routers

import (
	"net/http"
	"social/internal/handlers"
	"social/internal/middleware"
	"social/internal/utils"
	"time"
)

func Auhentication(servermux *http.ServeMux, handler *handlers.Handlers ,midware *middleware.Miiiddeware) {
	neewratelimiter := utils.NewRateLimiter(time.Minute)
	//alll auth routes
	servermux.Handle("/Me", neewratelimiter.ApplyMiddleware("authintication", (http.HandlerFunc(handler.ME))))
	servermux.Handle("/Signup", neewratelimiter.ApplyMiddleware("authintication", (http.HandlerFunc(handler.SignUp))))
	servermux.Handle("/Signin", neewratelimiter.ApplyMiddleware("authintication", (http.HandlerFunc(handler.Login))))
	servermux.Handle("/Logout", neewratelimiter.ApplyMiddleware("authintication", (http.HandlerFunc(handler.Logout))))
}
