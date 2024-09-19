package storage

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DBstorage struct {
	Pool *pgxpool.Pool
}
// Создание нового пула соединений
func NewDB(pool *pgxpool.Pool) (*DBstorage, error) {
	if pool == nil {
		return nil, fmt.Errorf("invalid pool: cannot be nil")
	}
	return &DBstorage{
		Pool: pool,
	}, nil
}

// Закрытие пула соединений
func (db *DBstorage) Close() {
	db.Pool.Close()
}
