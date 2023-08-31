package dto

import "time"

type ConsumerRequest struct {
	ConsumerId int            `json:"consumer_id"`
	Segments   []DataSegments `json:"segments"`
}

type ConsumerRequestDelete struct {
	ConsumerId int        `json:"consumer_id"`
	Segments   []Segments `json:"segments"`
}

type Segments struct {
	SegmentName string `json:"segment_name"`
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
