package entity

type Users struct {
	Id       int    `db:"id"`
	Username string `db:"username"`
	Password string `db:"password"`
}
