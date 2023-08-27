package entity

type Action string

type Actions struct {
	Id   int    `db:"id"`
	Name string `db:"name"`
}
