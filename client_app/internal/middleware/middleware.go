package middleware

import "social/internal/services"
type Miiiddeware struct{
	Serviiice *services.Service
}
func Neewmidd(service *services.Service)*Miiiddeware{
	return &Miiiddeware{
		Serviiice: service,
	}
}