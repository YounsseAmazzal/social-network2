package handlers

import "social/internal/services"

type Handlers struct {
	services   *services.Service
	wsmaanager *Wshandler
}

func NeewHandler(seerivcces *services.Service, wsmaanager *Wshandler) *Handlers {
	return &Handlers{
		services:   seerivcces,
		wsmaanager: wsmaanager,
	}
}
