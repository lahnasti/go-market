package storage

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/lahnasti/go-market/internal/models"
	"github.com/lahnasti/go-market/internal/repository"
)

type RepoProduct struct {
	db *DBstorage
}

func NewRepoProduct(db *DBstorage) repository.ProductRepository {
	return &RepoProduct{db: db}
}

func (r *RepoProduct) GetAllProducts() ([]models.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := r.db.Pool.Query(ctx, "SELECT uid, name, description, price, delete FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products []models.Product
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.UID, &product.Name, &product.Description, &product.Price, &product.Delete); err != nil {
			return nil, err
		}
		product.Name = strings.TrimSpace(product.Name)
		product.Description = strings.TrimSpace(product.Description)
		products = append(products, product)
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

func (r *RepoProduct) UpdateProduct(uid int, product models.Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := r.db.Pool.Exec(ctx, "UPDATE products SET name=$1, description=$2, price=$3 WHERE uid=$4", product.Name, product.Description, product.Price, uid)
	if err != nil {
		return fmt.Errorf("update user failed: %w", err)
	}
	return nil
}

func (r *RepoProduct) SetDeleteStatus(bid int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := r.db.Pool.Exec(ctx, "UPDATE books SET delete = true WHERE bid = $1"); err != nil {
		return fmt.Errorf("update delete status failed: %w", err)
	}
	return nil
}

func (r *RepoProduct) DeleteProduct() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	tx, err := r.db.Pool.Begin(ctx)
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
