package service

import (
	"context"
	"dynamic-user-segmentation/internal/repository"
	"dynamic-user-segmentation/internal/repository/respository_errors"
	"dynamic-user-segmentation/internal/service/dto"
	e "dynamic-user-segmentation/pkg/util/errors"
)

type SegmentsService struct {
	segments repository.Segments
}

func NewSegmentsService(segmentsRepository repository.Segments) *SegmentsService {
	return &SegmentsService{segments: segmentsRepository}
}

func (ss *SegmentsService) CreateSegment(ctx context.Context, segment dto.SegmentsRequest) (int, error) {
	return ss.segments.CreateSegment(ctx, segment.Name, segment.Percent)
}
func (ss *SegmentsService) DeleteSegment(ctx context.Context, segment dto.SegmentsRequest) error {
	var err error
	defer func() {
		err = e.WrapIfErr("Segments service: ", err)
	}()

	id, err := ss.segments.GetIdBySegment(ctx, segment.Name)
	if err != nil {
		return respository_errors.ErrNotFound
	}
	return ss.segments.DeleteSegment(ctx, id)
}
func (ss *SegmentsService) UpdateSegment(ctx context.Context, segment dto.SegmentsRequest) error {
	id, err := ss.segments.GetIdBySegment(ctx, segment.Name)
	if err != nil {
		return respository_errors.ErrNotFound
	}
	return ss.segments.UpdateSegment(ctx, id, segment.Percent)
}
func (ss *SegmentsService) GetAllSegments(ctx context.Context) ([]dto.SegmentsResponse, error) {
	res, err := ss.segments.GetAllSegments(ctx)
	if err != nil {
		return nil, e.Wrap("can't get all segments: ", err)
	}

	segments := make([]dto.SegmentsResponse, 0)
	for _, value := range res {
		var segment dto.SegmentsResponse
		segment.Id = value.Id
		segment.Name = value.Name
		segment.Percent = value.Percent
		segments = append(segments, segment)
	}
	return segments, nil
}
