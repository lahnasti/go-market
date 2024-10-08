package repository

import "github.com/lahnasti/go-market/internal/models"

type Repository interface {
	ProductRepository
	PurchaseRepository
	UserRepository
}

type PurchaseRepository interface {
	MakePurchase(models.Purchase) (int, error)
	GetUserPurchases(int) ([]models.Purchase, error)
	GetProductPurchases(int) ([]models.Purchase, error)
}

type ProductRepository interface {
	GetAllProducts() ([]models.Product, error)
	GetProductByID(int) (models.Product, error)
	AddProduct(models.Product) (int, error)
	UpdateProduct(int, models.Product) (int, error)
	DeleteProducts() error
	SetDeleteStatus(int) error
	IsProductUnique(string)(bool, error)
}

type UserRepository interface {
	GetUserProfile(int) (models.User, error)
	RegisterUser(models.User) (int, error)
	LoginUser(string, string) (int, error)
	IsUsernameUnique(string)(bool, error)
}