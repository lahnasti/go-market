package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/lahnasti/go-market/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (db *DBstorage) GetUserProfile(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	row := db.Pool.QueryRow(ctx, "SELECT * FROM users WHERE id=$1", id)
	var user models.User
	//Нужно ли пароль выводить?
	if err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Email); err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (db *DBstorage) RegisterUser(user models.User) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	row := db.Pool.QueryRow(ctx, "INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id", user.Username, user.Email, user.Password)
	var ID int
	if err := row.Scan(&ID); err != nil {
		return -1, err
	}
	return ID, nil
}

func (db *DBstorage) LoginUser(username string, password string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var userID int
	var hashedPassword string
	row := db.Pool.QueryRow(ctx, "SELECT id, password FROM users WHERE username=$1", username)
	if err := row.Scan(&userID, &hashedPassword); err != nil {
		if err == sql.ErrNoRows {
			return -1, fmt.Errorf("user not found")
		}
		return -1, err
	}
	// Проверка введенного пароля
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		// Пароль неверный
		return -1, fmt.Errorf("invalid password")
	}

	// Успешная авторизация
	return userID, nil
}

func (db *DBstorage) IsUsernameUnique(username string)(bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var count int
	row := db.Pool.QueryRow(ctx, "SELECT COUNT(*) FROM users WHERE username=$1", username)
    if err := row.Scan(&count); err!= nil {
        return false, fmt.Errorf("failed to check username existence: %w", err)
    }
	return count == 0, nil
}