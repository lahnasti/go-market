package models

// UserCheckMessage представляет структуру сообщения для проверки пользователя
type UserCheckMessage struct {
	UserID int `json:"userID"`
}

// UserCheckResponse представляет структуру ответа на проверку пользователя
type UserCheckResponse struct {
	Valid bool        `json:"valid"`
	User User `json:"user,omitempty"`
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