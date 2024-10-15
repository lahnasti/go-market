package models

import "time"

type Purchase struct {
	UID        int       `json:"uid"`
	UserID    int       `json:"userID" validate:"required"`
	ProductID int       `json:"productID" validate:"required"`
	Quantity  int       `json:"quantity" validate:"required"`
	PurchaseDate  time.Time `json:"purchase_date"`
}
