package repository

type Repository interface {
	ProductRepository
	PurchaseRepository
	UserRepository
}