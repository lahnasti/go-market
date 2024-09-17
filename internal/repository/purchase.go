package repository

import "github.com/lahnasti/go-market/internal/models"

type PurchaseRepository interface {
	MakePurchase(models.Purchase)(int, error)
	GetUserPurchases(int)([]models.Purchase, error)
	GetProductPurchases(int)([]models.Purchase, error)
}