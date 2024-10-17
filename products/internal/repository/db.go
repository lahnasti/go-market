package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

// EnsureAuthDatabaseExists проверяет наличие базы данных для сервиса auth и создаёт её при необходимости.
func EnsureMarketDatabaseExists(conn *pgxpool.Conn) error {
	const dbName = "market"
	return EnsureDatabaseExists(conn, dbName)
}

// Общая функция для проверки и создания базы данных
func EnsureDatabaseExists(conn *pgxpool.Conn, dbName string) error {
	// Проверяем, существует ли база данных
	var exists bool
	err := conn.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)", dbName).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if database exists: %w", err)
	}

	// Если база данных не существует, создаём её
	if !exists {
		_, err = conn.Exec(context.Background(), fmt.Sprintf("CREATE DATABASE %s", dbName))
		if err != nil {
			return fmt.Errorf("failed to create database: %w", err)
		}
		log.Printf("Database %s created successfully", dbName)
	}

	return nil
}
