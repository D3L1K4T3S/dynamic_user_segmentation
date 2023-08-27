package entity

import "time"

type ConsumersSegments struct {
	Id        int       `db:"id"`
	SegmentId int       `db:"segment_id"`
	TTL       time.Time `db:"ttl"`
}
