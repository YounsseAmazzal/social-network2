package repositories

import "social/internal/models"

type Repository struct {
	DB *models.DB
}
func Neewwrepo(db *models.DB)*Repository{
	return &Repository{
		DB: db,
	}
}