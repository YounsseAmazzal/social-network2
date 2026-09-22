package handlers

import "social/internal/services"
type Wshandler struct{
	seervice *services.Service
}
func Neeewsmanger(sevice *services.Service)*Wshandler{
	return &Wshandler{
		seervice: sevice,
	}
}