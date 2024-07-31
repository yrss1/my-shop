package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
)

type ProductRepository struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) List(ctx context.Context) (dest []product.Entity, err error) {
	query := `
			SELECT id, name, description, price, category, quantity 
			FROM products
			ORDER BY id`

	err = r.db.SelectContext(ctx, &dest, query)

	return
}

func (r *ProductRepository) Add(ctx context.Context, data product.Entity) (id string, err error) {
	query := `
		INSERT INTO products (name, description, price, category, quantity) 
		VALUES ($1, $2, $3, $4, $5) 
		RETURNING id`

	args := []any{data.Name, data.Description, data.Price, data.Category, data.Quantity}

	if err = r.db.QueryRowContext(ctx, query, args...).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = store.ErrorNotFound
		}
	}

	return
}

func (r *ProductRepository) Get(ctx context.Context, id string) (dest product.Entity, err error) {
	query := `
			SELECT id, name, description, price, category, quantity 
			FROM products
			WHERE id=$1`

	args := []any{id}

	if err = r.db.GetContext(ctx, &dest, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = store.ErrorNotFound
		}
	}

	return
}

func (r *ProductRepository) Update(ctx context.Context, id string, data product.Entity) (err error) {
	sets, args := r.prepareArgs(data)

	if len(args) > 0 {
		args = append(args, id)
		sets = append(sets, "updated_at=CURRENT_TIMESTAMP")

		query := fmt.Sprintf("UPDATE products SET %s WHERE id=$%d RETURNING id", strings.Join(sets, ", "), len(args))

		if err = r.db.QueryRowContext(ctx, query, args...).Scan(&id); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				err = store.ErrorNotFound
			}
		}
	}

	return
}

func (r *ProductRepository) prepareArgs(data product.Entity) (sets []string, args []any) {
	if data.Name != nil {
		args = append(args, data.Name)
		sets = append(sets, fmt.Sprintf("name=$%d", len(args)))
	}

	if data.Description != nil {
		args = append(args, data.Description)
		sets = append(sets, fmt.Sprintf("description=$%d", len(args)))
	}

	if data.Price != nil {
		args = append(args, data.Price)
		sets = append(sets, fmt.Sprintf("price=$%d", len(args)))
	}

	if data.Category != nil {
		args = append(args, data.Category)
		sets = append(sets, fmt.Sprintf("category=$%d", len(args)))
	}

	if data.Quantity != nil {
		args = append(args, data.Quantity)
		sets = append(sets, fmt.Sprintf("quantity=$%d", len(args)))
	}

	return
}

func (r *ProductRepository) Delete(ctx context.Context, id string) (err error) {
	query := `
		DELETE FROM products
		WHERE id=$1
		RETURNING id`

	args := []any{id}

	if err = r.db.QueryRowContext(ctx, query, args...).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = store.ErrorNotFound
		}
	}

	return
}

func (r *ProductRepository) Search(ctx context.Context, data product.Entity) (dest []product.Entity, err error) {
	query := "SELECT id, name, description, price, category, quantity FROM products WHERE 1=1"

	sets, args := r.prepareArgs(data)
	if len(sets) > 0 {
		query += " AND " + strings.Join(sets, " AND ")
	}

	err = r.db.SelectContext(ctx, &dest, query, args...)

	return
}
