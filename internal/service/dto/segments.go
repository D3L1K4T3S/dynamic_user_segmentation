package dbo

type SegmentsRequest struct {
	Name    string
	Percent float64
}

type SegmentsResponse struct {
	Id      int
	Name    string
	Percent float64
}
