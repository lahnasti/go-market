package models

import "time"

// UserCheckMessage представляет структуру сообщения для проверки пользователя
type UserCheckMessage struct {
	UserID int `json:"userID"`
}

// UserCheckResponse представляет структуру ответа на проверку пользователя
type UserCheckResponse struct {
	Valid bool `json:"valid"`
	User  User `json:"user,omitempty"`
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=8"`
	Email    string `json:"email" validate:"required,email"`
}

type Credentials struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Product struct {
	UID         int     `json:"uid"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Price       float64 `json:"price" validate:"required,min=1"`
	Delete      bool    `json:"delete"`
	Quantity    int     `json:"quantity" validate:"required,min=1"`
}

type Purchase struct {
	UID          int       `json:"uid"`
	UserID       int       `json:"userID" validate:"required"`
	ProductID    int       `json:"productID" validate:"required"`
	Quantity     int       `json:"quantity" validate:"required"`
	PurchaseDate time.Time `json:"purchase_date"`
}
