package handlers

import "social/internal/services"
type handlers struct{
	services *services.Service
	wsmaanager *wshandler
}
func  NeewHandler(seerivcces *services.Service,wsmaanager *wshandler) *handlers {
	return &handlers{
		services: seerivcces,
		wsmaanager: wsmaanager,
	}
}