package user

type Entity struct {
	ID       string  `db:"id"`
	Name     *string `db:"name"`
	Email    *string `db:"email"`
	Password *string `db:"password"`
	Address  *string `db:"address"`
	Role     *string `db:"role"`
}
