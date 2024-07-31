package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
)

type PaymentRepository struct {
	db *sqlx.DB
}

func NewPaymentRepository(db *sqlx.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) List(ctx context.Context) (dest []payment.Entity, err error) {
	query := `
			SELECT id, user_id, order_id, amount, status
			FROM payments
			ORDER BY id`

	err = r.db.SelectContext(ctx, &dest, query)

	return
}

func (r *PaymentRepository) Add(ctx context.Context, data payment.Entity) (id string, err error) {
	query := `
		INSERT INTO payments (user_id, order_id, amount, status) 
		VALUES ($1, $2, $3, $4) 
		RETURNING id`

	args := []any{data.UserID, data.OrderID, data.Amount, data.Status}

	if err = r.db.QueryRowContext(ctx, query, args...).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = store.ErrorNotFound
		}
	}

	return
}

func (r *PaymentRepository) Get(ctx context.Context, id string) (dest payment.Entity, err error) {
	query := `
		SELECT id, user_id, order_id, amount, status 
		FROM payments 
		WHERE id=$1`

	args := []any{id}

	if err = r.db.GetContext(ctx, &dest, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			err = store.ErrorNotFound
		}
	}

	return
}

func (r *PaymentRepository) Update(ctx context.Context, id string, data payment.Entity) (err error) {
	sets, args := r.prepareArgs(data)

	if len(args) > 0 {
		args = append(args, id)
		sets = append(sets, "updated_at=CURRENT_TIMESTAMP")

		query := fmt.Sprintf("UPDATE payments SET %s WHERE id=$%d RETURNING id", strings.Join(sets, ", "), len(args))

		if err = r.db.QueryRowContext(ctx, query, args...).Scan(&id); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				err = store.ErrorNotFound
			}
		}
	}

	return
}

func (r *PaymentRepository) prepareArgs(data payment.Entity) (sets []string, args []any) {
	if data.UserID != nil {
		args = append(args, data.UserID)
		sets = append(sets, fmt.Sprintf("user_id=$%d", len(args)))
	}

	if data.OrderID != nil {
		args = append(args, data.OrderID)
		sets = append(sets, fmt.Sprintf("order_id=$%d", len(args)))
	}

	if data.Amount != nil {
		args = append(args, data.Amount)
		sets = append(sets, fmt.Sprintf("amount=$%d", len(args)))
	}

	if data.Status != nil {
		args = append(args, data.Status)
		sets = append(sets, fmt.Sprintf("status=$%d", len(args)))
	}

	return
}

func (r *PaymentRepository) Delete(ctx context.Context, id string) (err error) {
	query := `
		DELETE FROM payments
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

func (r *PaymentRepository) Search(ctx context.Context, data payment.Entity) (dest []payment.Entity, err error) {
	query := "SELECT id, user_id, order_id, amount, status FROM payments WHERE 1=1"

	sets, args := r.prepareArgs(data)
	if len(sets) > 0 {
		query += " AND " + strings.Join(sets, " AND ")
	}

	err = r.db.SelectContext(ctx, &dest, query, args...)

	return
}
