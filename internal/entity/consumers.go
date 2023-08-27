package entity

type Consumers struct {
	Id         int `db:"id"`
	ConsumerId int `db:"consumer_id"`
	SegmentId  int `db:"segment_id"`
}
