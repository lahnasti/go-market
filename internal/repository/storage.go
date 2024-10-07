package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBstorage struct {
	Pool *pgxpool.Pool
}

// Создание нового пула соединений
func NewDB(pool *pgxpool.Pool) (*DBstorage, error) {
	return &DBstorage{
		Pool: pool,
	}, nil
}

// Закрытие пула соединений
func (db *DBstorage) Close() {
	db.Pool.Close()
}
