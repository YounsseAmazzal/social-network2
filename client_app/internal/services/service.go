package services

import "social/internal/repositories"

type Service struct {
	Reepo *repositories.Repository
}

func NewService(reepo *repositories.Repository) *Service {
	return &Service{Reepo: reepo}
}
