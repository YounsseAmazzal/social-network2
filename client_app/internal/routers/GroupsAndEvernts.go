package routers

import (
	"net/http"
	"social/internal/handlers"
	"social/internal/middleware"
	"social/internal/utils"
	"time"
)

func Groupsandevents(servermux *http.ServeMux, handler *handlers.Handlers, midware *middleware.Miiiddeware) {
	neewratelimiter := utils.NewRateLimiter(time.Minute)
	//alll groups  routes
	servermux.Handle("/CreateGroup", neewratelimiter.ApplyMiddleware("normalapi", (http.HandlerFunc(handler.CreateGroup))))
	servermux.Handle("/DiscoverGroups", neewratelimiter.ApplyMiddleware("normalapi", (http.HandlerFunc(handler.DiscooverGroups))))
	servermux.Handle("/JoinedGroups", neewratelimiter.ApplyMiddleware("normalapi", (http.HandlerFunc(handler.JoinedGroups))))
	servermux.Handle("/SeeeTheGroup", neewratelimiter.ApplyMiddleware("normalapi", (http.HandlerFunc(handler.SeeeTheGroup))))
	servermux.Handle("/GrouupMembers", neewratelimiter.ApplyMiddleware("normalapi", (http.HandlerFunc(handler.GrouupMembers))))
	servermux.Handle("/CreatePostGroup", neewratelimiter.ApplyMiddleware("normalapi", (http.HandlerFunc(handler.CreatePostGroup))))
	servermux.Handle("/GetGroupPosts", neewratelimiter.ApplyMiddleware("normalapi", (http.HandlerFunc(handler.GroupPosts))))
	// servermux.Handle("/CreateGPComment", neewratelimiter.ApplyMiddleware("normalapi", (http.HandlerFunc(handler.CreateGPComment))))
	servermux.Handle("/GroupCreateEvent", neewratelimiter.ApplyMiddleware("normalapi", (http.HandlerFunc(handler.GroupCreateEvent))))
	servermux.Handle("/GroupEvents", neewratelimiter.ApplyMiddleware("normalapi", (http.HandlerFunc(handler.GroupEvents))))
	servermux.Handle("/GroupEventParticipate", neewratelimiter.ApplyMiddleware("normalapi", (http.HandlerFunc(handler.GroupEventParrticipate))))

	// servermux.Handle("/Postlike", neewratelimiter.ApplyMiddleware("normalapi", (http.HandlerFunc(handler.LikePost))))
	// servermux.Handle("/Postdislike", neewratelimiter.ApplyMiddleware("normalapi", (http.HandlerFunc(handler.DislikePost))))
	// servermux.Handle("/GetComment", neewratelimiter.ApplyMiddleware("normalapi", (http.HandlerFunc(handler.GetComment))))
	// servermux.Handle("/AllPostComment", neewratelimiter.ApplyMiddleware("normalapi", (http.HandlerFunc(handler.GetComments))))
	// servermux.Handle("/RatePost", neewratelimiter.ApplyMiddleware("normalapi", (http.HandlerFunc(handler.AllComments))))
	// servermux.Handle("/Ratecomment", neewratelimiter.ApplyMiddleware("normalapi", (http.HandlerFunc(handler.AllComments))))
}
