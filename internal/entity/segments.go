package entity

type Segments struct {
	Id      int     `db:"id"`
	Name    string  `db:"name"`
	Percent float64 `db:"percent"`
}
