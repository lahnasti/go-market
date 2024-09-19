package storage

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/lahnasti/go-market/internal/models"
)
func (db *DBstorage) GetAllProducts() ([]models.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := db.Pool.Query(ctx, "SELECT uid, name, description, price, delete, quantity FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products []models.Product
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.UID, &product.Name, &product.Description, &product.Price, &product.Delete, &product.Quantity); err != nil {
			return nil, err
		}
		product.Name = strings.TrimSpace(product.Name)
		product.Description = strings.TrimSpace(product.Description)
		products = append(products, product)
	}
	return products, nil
}

func (db *DBstorage) GetProductByID(uid int) (*models.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	row := db.Pool.QueryRow(ctx, "SELECT * FROM product WHERE uid=$1", uid)
	var product models.Product
	if err := row.Scan(&product.UID, &product.Name, &product.Description, &product.Price, &product.Delete, &product.Quantity); err != nil {
		return nil, err
	}
	product.Name = strings.TrimSpace(product.Name)
	product.Description = strings.TrimSpace(product.Description)
	return &product, nil
}

func (db *DBstorage) AddProduct(product models.Product) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	row := db.Pool.QueryRow(ctx, "INSERT INTO products (name, description, price, quantity) VALUES ($1, $2, $3, $4) RETURNING uid", product.Name, product.Description, product.Price, product.Quantity)
	var UID int
	if err := row.Scan(&UID); err != nil {
		return -1, err
	}
	return UID, nil
}

func (db *DBstorage) UpdateProduct(uid int, product models.Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := db.Pool.Exec(ctx, "UPDATE products SET name=$1, description=$2, price=$3, quantity=$4 WHERE uid=$5", product.Name, product.Description, product.Price, product.Quantity, uid)
	if err != nil {
		return fmt.Errorf("update user failed: %w", err)
	}
	return nil
}

func (db *DBstorage) SetDeleteStatus(uid int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := db.Pool.Exec(ctx, "UPDATE products SET delete = true WHERE uid = $1", uid); err != nil {
		return fmt.Errorf("update delete status failed: %w", err)
	}
	return nil
}

func (db *DBstorage) DeleteProducts() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("create transaction failed: %w", err)
	}
	defer tx.Rollback(ctx)
	if _, err := tx.Prepare(ctx, "delete product", "DELETE FROM products WHERE delete = true"); err != nil {
		return fmt.Errorf("create prepare sql str failed: %w", err)
	}
	if _, err := tx.Exec(ctx, "delete product"); err != nil {
		return fmt.Errorf("failed delete product: %w", err)
	}
	return tx.Commit(ctx)
}
