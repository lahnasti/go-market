package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/lahnasti/go-market/internal/models"
)

func (db *DBstorage) MakePurchase(purchase models.Purchase) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return -1, err
	}
	defer tx.Rollback(ctx)
	// Уменьшаем количество продукта в таблице products
	updateQuery := `UPDATE products SET quantity = quantity - $1 WHERE uid = $2 AND quantity >= $1`
	result, err := tx.Exec(ctx, updateQuery, purchase.Quantity, purchase.ProductID)
	if err != nil {
		return -1, err
	}
	// Проверяем, затронута ли строка (т.е. продукт в наличии в нужном количестве)
	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return -1, fmt.Errorf("not enough product quantity available or product does not exist")
	}

	insertQuery := `INSERT INTO purchases (user_id, product_id, quantity) VALUES ($1, $2, $3) RETURNING uid`
	row := tx.QueryRow(ctx, insertQuery, purchase.UserID, purchase.ProductID, purchase.Quantity)
	var UID int
	if err := row.Scan(&UID); err != nil {
		return -1, err
	}
	// Фиксируем транзакцию
	if err := tx.Commit(ctx); err != nil {
		return -1, err
	}
	return UID, nil
}

func (db *DBstorage) GetUserPurchases(userID int) ([]models.Purchase, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := db.Pool.Query(ctx, "SELECT * FROM purchases WHERE user_id=$1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var purchases []models.Purchase
	for rows.Next() {
		var purchase models.Purchase
		if err := rows.Scan(&purchase.UID, &purchase.UserID, &purchase.ProductID, &purchase.Quantity); err != nil {
			return nil, err
		}
		purchases = append(purchases, purchase)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return purchases, nil
}

func (db *DBstorage) GetProductPurchases(productID int) ([]models.Purchase, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := db.Pool.Query(ctx, "SELECT * FROM purchases WHERE product_id=$1", productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var purchases []models.Purchase
	for rows.Next() {
		var purchase models.Purchase
		if err := rows.Scan(&purchase.UID, &purchase.UserID, &purchase.ProductID, &purchase.Quantity); err != nil {
			return nil, err
		}
		purchases = append(purchases, purchase)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return purchases, nil
}
