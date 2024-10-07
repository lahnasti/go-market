package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/huandu/go-sqlbuilder"
	"github.com/lahnasti/go-market/internal/models"
)

func (db *DBstorage) GetAllProducts() ([]models.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sb := sqlbuilder.NewSelectBuilder()
	query, args := sb.Select("*").From("products").Build()

	rows, err := db.Pool.Query(ctx, query, args...)
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

func (db *DBstorage) GetProductByID(uid int) (models.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sb := sqlbuilder.NewSelectBuilder()
	query, args := sb.Select("*").From("products").Where(sb.Equal("uid", uid)).Build()

	row := db.Pool.QueryRow(ctx, query, args...)
	var product models.Product
	if err := row.Scan(&product.UID, &product.Name, &product.Description, &product.Price, &product.Delete, &product.Quantity); err != nil {
		return models.Product{}, err
	}
	product.Name = strings.TrimSpace(product.Name)
	product.Description = strings.TrimSpace(product.Description)
	return product, nil
}

func (db *DBstorage) AddProduct(product models.Product) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sb := sqlbuilder.NewInsertBuilder()
	query, args := sb.InsertInto("products").Cols("name", "description", "price", "quantity").
					Values(product.Name, product.Description, product.Price, product.Quantity).
					BuildWithFlavor(sqlbuilder.PostgreSQL)
					query += "RETURNING uid"

	var UID int
	err := db.Pool.QueryRow(ctx, query, args...).Scan(&UID)
	if err != nil {
		return -1, fmt.Errorf("failed to insert product: %w", err)
	}
	return UID, nil
}

func (db *DBstorage) UpdateProduct(uid int, product models.Product) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	sb := sqlbuilder.NewUpdateBuilder()
	query, args := sb.Update("products").
					Set(
						sb.Assign("name", product.Name),
					    sb.Assign("description", product.Description),
                        sb.Assign("price", product.Price),
                        sb.Assign("quantity", product.Quantity),
                    ).
					Where(sb.Equal("uid", product.UID)).
					Build()
					query += "RETURNING uid"

	var UID int
	err := db.Pool.QueryRow(ctx, query, args...).Scan(&UID)
	if err != nil {
		return -1, fmt.Errorf("update user failed: %w", err)
	}
	return UID, nil
}

func (db *DBstorage) SetDeleteStatus(uid int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sb := sqlbuilder.NewUpdateBuilder()
	query, args := sb.Update("products").
					Set("delete", "true").Where(sb.Equal("uid", uid)).Build()

	_, err := db.Pool.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to set delete status: %w", err) 
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
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	sb := sqlbuilder.NewDeleteBuilder()
	query, args := sb.DeleteFrom("products").
					Where(sb.Equal("delete", true)).
					Build()

	if _, err := tx.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("create prepare sql str failed: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed delete product: %w", err)
	}
	return nil
}
