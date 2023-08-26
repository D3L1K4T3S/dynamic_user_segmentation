package entity

type Action string

const (
	Add    Action = "Add"
	Create Action = "Create"
	Delete Action = "Delete"
	Update Action = "Update"
)

type Actions struct {
	Id   int    `db:"id"`
	Name string `db:"name"`
}
