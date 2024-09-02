package storage

import (
	"context"
	"strings"
	"time"

	"github.com/lahnasti/go-market/internal/models"
	"github.com/lahnasti/go-market/internal/storage"
)

type RepoProduct struct {
	db *storage.DBstorage
}

func (r *RepoProduct) GetAllProducts() ([]models.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := r.db.Pool.Query(ctx, "SELECT * FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products []models.Product
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.UID, &product.Name, &product.Description, &product.Price); err != nil {
			return nil, err
		}
		product.Name = strings.TrimSpace(product.Name)
		product.Description = strings.TrimSpace(product.Description)
	}
	return products, nil
}

func (r *RepoProduct) AddProduct(product models.Product) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	row := r.db.Pool.QueryRow(ctx, "INSERT INTO products (name, description, price) VALUES ($1, $2, $3) RETURNING uid", product.Name, product.Description, product.Price)
	var UID int
	if err := row.Scan(&UID); err != nil {
		return -1, err
	}
	return UID, nil
}

func (r *RepoProduct) UpdateProduct(int, product models.Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := r.db.Pool.Exec(ctx, "UPDATE products SET name=$1, description=$2, price=$3 WHERE uid=$4", product.Name, product.Description, product.Price, product.UID)
	if err != nil {
		return nil
	}
	return nil
}

func (r *RepoProduct) DeleteProduct(uid int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := r.db.Pool.Exec(ctx, "DELETE FROM product WHERE uid=$1", uid)
	if err != nil {
		return err
	}
	return nil
}
