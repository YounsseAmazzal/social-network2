package requests

import (
	"errors"
	"net/http"
)

func ValidFollowReq(req *http.Request ,userviewerif string) (string, error) {
	TheFollowed := req.URL.Query().Get("Userfollowinng_id")
	if TheFollowed == "" {
		return "", errors.New("The taget follow not found")
	}
	if TheFollowed==userviewerif{
		return "", errors.New("You Couldn't follow your self hh")
	}
	return TheFollowed,nil
}
func ValidUnfollow(req *http.Request ,userviewerif string) (string, error){
	TheFollowed := req.URL.Query().Get("UserUnfollow_id")
	if TheFollowed == "" {
		return "", errors.New("The taget Unfollow not found")
	}
	if TheFollowed==userviewerif{
		return "", errors.New("You Couldn't Unfollow your self hh")
	}
	return TheFollowed,nil
}
func Accepting(req *http.Request ,userviewerif string) (string, error){
	TheFollowed := req.URL.Query().Get("FollowAccet_id")
	if TheFollowed == "" {
		return "", errors.New("The taget To accet follow not found")
	}
	if TheFollowed==userviewerif{
		return "", errors.New("You Couldn't Accept your self hh")
	}
	return TheFollowed,nil
}
func Declineinng(req *http.Request ,userviewerif string) (string, error){
	TheFollowed := req.URL.Query().Get("FollowDecline_id")
	if TheFollowed == "" {
		return "", errors.New("The taget To declin follow not found")
	}
	if TheFollowed==userviewerif{
		return "", errors.New("You Couldn't Decline your self hh")
	}
	return TheFollowed,nil
}