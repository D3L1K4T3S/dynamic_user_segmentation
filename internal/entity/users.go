package entity

type Users struct {
	Id       int `db:"id"`
	Username int `db:"username"`
	Password int `db:"password"`
}
