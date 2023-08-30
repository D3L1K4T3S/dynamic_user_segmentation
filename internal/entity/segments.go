package entity

import "time"

type Segments struct {
	Id         int       `db:"id"`
	Name       string    `db:"name"`
	Percent    float64   `db:"percent"`
	ModifiedAt time.Time `db:"modified_at"`
}
