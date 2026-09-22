package handlers

import "social/internal/services"
type wshandler struct{
	seervice *services.Service
}
func Neeewsmanger(sevice *services.Service)*wshandler{
	return &wshandler{
		seervice: sevice,
	}
}