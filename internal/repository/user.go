package repository

import "github.com/lahnasti/go-market/internal/models"

type UserRepository interface {
	GetUserProfile(int) (models.User, error)
	RegisterUser(models.User) (int, error)
	LoginUser(string, string) (int, error)
}