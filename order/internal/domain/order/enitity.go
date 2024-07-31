package order

type Entity struct {
	ID         string   `db:"id"`
	UserID     *string  `db:"user_id"`
	Products   []string `db:"products"`
	TotalPrice *float64 `db:"total_price"`
	Status     *string  `db:"status"`
}
