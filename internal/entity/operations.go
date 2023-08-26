package entity

import "time"

type Operations struct {
	Id        int       `db:"id"`
	UserId    int       `db:"user_id"`
	SegmentId int       `db:"segment_id"`
	ActionId  int       `db:"operation_id"`
	Created   time.Time `db:"created_at"`
}
