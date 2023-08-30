package dto

import "time"

type ConsumerRequest struct {
	ConsumerId int            `json:"consumer_id"`
	Segments   []DataSegments `json:"segments"`
}

type DataSegments struct {
	SegmentName string    `json:"segment_name"`
	TTL         time.Time `json:"ttl"`
}

type ConsumerResponse struct {
	ConsumerId   int      `json:"consumer_id"`
	SegmentsName []string `json:"segments_name"`
}

type ConsumerId struct {
	ConsumerId int `json:"consumer_id"`
}
