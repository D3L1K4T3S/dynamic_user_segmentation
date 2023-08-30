package dto

import "time"

type OperationsRequest struct {
	ConsumerId int       `json:"consumer_id"`
	StartAt    time.Time `json:"start_at"`
	EndAt      time.Time `json:"end_at"`
}

type OperationsResponse struct {
	ConsumerId     int              `json:"consumer_id"`
	OperationsData []OperationsData `json:"operations_data"`
}

type OperationsData struct {
	SegmentName string    `json:"segment_name"`
	ActionName  string    `json:"action_name"`
	Date        time.Time `json:"date"`
}
