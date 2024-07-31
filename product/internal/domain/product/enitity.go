package product

type Entity struct {
	ID          string   `db:"id"`
	Name        *string  `db:"name"`
	Description *string  `db:"description"`
	Price       *float64 `db:"price"`
	Category    *string  `db:"category"`
	Quantity    *int     `db:"quantity"`
}
