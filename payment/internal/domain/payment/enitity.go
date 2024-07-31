package payment

type Entity struct {
	ID      string  `db:"id"`
	UserID  *string `db:"user_id"`
	OrderID *string `db:"order_id"`
	Amount  *string `db:"amount"`
	Status  *string `db:"status"`
}
