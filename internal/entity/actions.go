package entity

type Actions struct {
	Id   int    `db:"id"`
	Name string `db:"name"`
}

const (
	ActionTypeCreate = "create"
	ActionTypeAdd    = "add"
	ActionTypeDelete = "delete"
	ActionTypeUpdate = "update"
)
