package storage

type Employee struct {
	ID         int     `db:"id"`
	FirstName  string  `db:"first_name"`
	LastName   string  `db:"last_name"`
	Recipient  string  `db:"recipient"`
	Wage       float64 `db:"wage"`
	Department string  `db:"department"`
}
