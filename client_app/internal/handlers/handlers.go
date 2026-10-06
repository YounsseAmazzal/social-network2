package handlers

import "social/internal/services"

type Handlers struct {
	services   *services.Service
	wsmaanager *WSManager
}

func NeewHandler(seerivcces *services.Service) *Handlers {
	return &Handlers{
		services:   seerivcces,
	}
}
