package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
)

type OrderRepository struct {
	db *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) List(ctx context.Context) (dest []order.Entity, err error) {
	query := `
		SELECT id, user_id, total_price, status 
		FROM orders
		ORDER BY id`

	var orders []order.Entity
	err = r.db.SelectContext(ctx, &orders, query)
	if err != nil {
		return nil, err
	}

	for i := range orders {
		products, err := r.getProductsByOrderID(ctx, orders[i].ID)
		if err != nil {
			return nil, err
		}
		orders[i].Products = products
	}

	return orders, nil
}

func (r *OrderRepository) Add(ctx context.Context, data order.Entity) (id string, err error) {
	query := `
		INSERT INTO orders (user_id, total_price, status) 
		VALUES ($1, $2, $3) 
		RETURNING id`

	if err = r.db.QueryRowContext(ctx, query, data.UserID, data.TotalPrice, data.Status).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = store.ErrorNotFound
		}
	}

	productQuery := `
		INSERT INTO order_product (order_id, product_id) 
		VALUES ($1, $2)`

	for _, productID := range data.Products {
		if _, err = r.db.ExecContext(ctx, productQuery, id, productID); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				err = store.ErrorNotFound
			}
		}
	}
	return
}

func (r *OrderRepository) Get(ctx context.Context, id string) (dest order.Entity, err error) {
	query := `
		SELECT id, user_id, total_price, status 
		FROM orders
		WHERE id=$1`

	if err = r.db.GetContext(ctx, &dest, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = store.ErrorNotFound
			return
		}
		return
	}

	products, err := r.getProductsByOrderID(ctx, dest.ID)
	if err != nil {
		return
	}
	dest.Products = products

	return
}

func (r *OrderRepository) Update(ctx context.Context, id string, data order.Entity) (err error) {
	sets, args := r.prepareArgs(data)

	if len(args) > 0 {
		args = append(args, id)
		sets = append(sets, "updated_at=CURRENT_TIMESTAMP")

		query := fmt.Sprintf("UPDATE orders SET %s WHERE id=$%d RETURNING id", strings.Join(sets, ", "), len(args))

		if err = r.db.QueryRowContext(ctx, query, args...).Scan(&id); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				err = store.ErrorNotFound
				return
			}
		}
	}

	if data.Products != nil {
		if err = r.updateOrderProducts(ctx, id, data.Products); err != nil {
			return
		}
	}

	return
}

func (r *OrderRepository) Delete(ctx context.Context, id string) (err error) {
	query := `
		DELETE FROM orders
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

func (r *OrderRepository) Search(ctx context.Context, data order.Entity) (dest []order.Entity, err error) {
	query := "SELECT id, user_id, total_price, status FROM orders WHERE 1=1"

	sets, args := r.prepareArgs(data)
	if len(sets) > 0 {
		query += " AND " + strings.Join(sets, " AND ")
	}
	err = r.db.SelectContext(ctx, &dest, query, args...)
	if err != nil {
		return
	}

	for i := range dest {
		products, err := r.getProductsByOrderID(ctx, dest[i].ID)
		if err != nil {
			return nil, err
		}
		dest[i].Products = products
	}

	return
}

func (r *OrderRepository) prepareArgs(data order.Entity) (sets []string, args []any) {
	if data.UserID != nil {
		args = append(args, data.UserID)
		sets = append(sets, fmt.Sprintf("user_id=$%d", len(args)))
	}

	if data.TotalPrice != nil {
		args = append(args, data.TotalPrice)
		sets = append(sets, fmt.Sprintf("total_price=$%d", len(args)))
	}

	if data.Status != nil {
		args = append(args, data.Status)
		sets = append(sets, fmt.Sprintf("status=$%d", len(args)))
	}

	return
}

func (r *OrderRepository) getProductsByOrderID(ctx context.Context, orderID string) (products []string, err error) {
	query := `
		SELECT p.id
		FROM products p
		JOIN order_product op ON p.id = op.product_id
		WHERE op.order_id = $1`

	err = r.db.SelectContext(ctx, &products, query, orderID)
	if err != nil {
		return
	}

	return
}

func (r *OrderRepository) updateOrderProducts(ctx context.Context, orderID string, products []string) (err error) {
	deleteQuery := "DELETE FROM order_product WHERE order_id=$1"
	if _, err = r.db.ExecContext(ctx, deleteQuery, orderID); err != nil {
		return
	}

	productQuery := "INSERT INTO order_product (order_id, product_id) VALUES ($1, $2)"
	for _, productID := range products {
		if _, err = r.db.ExecContext(ctx, productQuery, orderID, productID); err != nil {
			return
		}
	}

	return
}
