package dto

type SegmentsRequest struct {
	Name    string  `json:"name"`
	Percent float64 `json:"percent"`
}

type SegmentsResponse struct {
	Id      int     `json:"id"`
	Name    string  `json:"name"`
	Percent float64 `json:"percent"`
}

const Null = "Null"
