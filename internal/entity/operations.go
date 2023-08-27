package entity

import "time"

type Operations struct {
	Id         int       `db:"id"`
	ConsumerId int       `db:"consumer_id"`
	SegmentId  int       `db:"segment_id"`
	ActionId   int       `db:"action_id"`
	Created    time.Time `db:"created_at"`
}
