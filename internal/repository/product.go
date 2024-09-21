package repository

import "github.com/lahnasti/go-market/internal/models"

type ProductRepository interface {
	GetAllProducts() ([]models.Product, error)
	GetProductByID(int) (*models.Product, error)
	AddProduct(models.Product) (int, error)
	UpdateProduct(int, models.Product) error
	DeleteProducts() error
	SetDeleteStatus(int) error
}