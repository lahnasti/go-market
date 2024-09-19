package models

import "time"

type Purchase struct {
	ID        int       `json:"id"`
	UserID    int       `json:"userId" validate:"required"`
	ProductID int       `json:"productId" validate:"required"`
	Quantity  int       `json:"quantity" validate:"required"`
	Timestamp time.Time `json:"timestamp"`
}
