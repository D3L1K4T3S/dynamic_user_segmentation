package entity

import "time"

type ComplexOperations struct {
	ConsumerId  int       `db:"consumer_id"`
	SegmentName string    `db:"segment_name"`
	ActionName  string    `db:"action_name"`
	Created     time.Time `db:"created_at"`
}

type ComplexConsumerSegments struct {
	ConsumerId  int    `db:"consumer_id"`
	SegmentName string `db:"segment_name"`
}
