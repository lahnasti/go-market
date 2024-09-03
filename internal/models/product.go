package models

type Product struct {
	UID         int     `json:"uid"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Price       float64 `json:"price" validate:"required, gt=0"`
	Delete      bool    `json:"delete"`
	Quantity    int     `json:"quantity" validate:"required"`
}
