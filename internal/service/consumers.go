package service

import (
	"context"
	"dynamic-user-segmentation/internal/entity"
	"dynamic-user-segmentation/internal/repository"
	"dynamic-user-segmentation/internal/repository/respository_errors"
	"dynamic-user-segmentation/internal/service/dto"
	"dynamic-user-segmentation/pkg/util/count_percent"
	e "dynamic-user-segmentation/pkg/util/errors"
)

type ConsumersService struct {
	consumers         repository.Consumers
	segments          repository.Segments
	consumersSegments repository.ConsumersSegments
	operations        repository.Operations
	actions           repository.Actions
}

func NewConsumersService(consumers repository.Consumers, segments repository.Segments, consumersSegments repository.ConsumersSegments) *ConsumersService {
	return &ConsumersService{
		consumers:         consumers,
		segments:          segments,
		consumersSegments: consumersSegments,
	}
}

func (cs *ConsumersService) CreateConsumer(ctx context.Context, consumer dto.ConsumerRequest) ([]int, error) {
	var err error
	defer func() {
		err = e.WrapIfErr("Service consumer: ", err)
	}()

	check, err := cs.consumers.GetSegmentsById(ctx, consumer.ConsumerId)
	if len(check) > 0 {
		return nil, ErrUserAlreadyExists
	}

	actionId, err := cs.actions.GetIdByAction(ctx, entity.ActionTypeCreate)
	if err != nil {
		return nil, e.Wrap("can't get action: ", err)
	}

	var consSegmentId int
	var segmentId int
	var res []int
	var id int
	for _, value := range consumer.Segments {
		segmentId, err = cs.segments.GetIdBySegment(ctx, value.SegmentName)
		if err != nil {
			return nil, respository_errors.ErrNotFound
		}

		if value.TTL.IsZero() {
			consSegmentId, err = cs.consumersSegments.AddConsumerSegmentTTL(ctx, segmentId, value.TTL)
			if err != nil {
				return nil, err
			}
		} else {
			consSegmentId, err = cs.consumersSegments.AddConsumerSegment(ctx, segmentId)
			if err != nil {
				return nil, err
			}
		}

		id, err = cs.consumers.CreateConsumer(ctx, consSegmentId, consumer.ConsumerId)
		if err != nil {
			return nil, respository_errors.ErrCannotCreate
		}
		res = append(res, id)

		_, err = cs.operations.AddOperation(ctx, consumer.ConsumerId, consSegmentId, actionId)
		if err != nil {
			return nil, respository_errors.ErrCannotCreate
		}
	}

	extRes, err := cs.automaticAdd(ctx, consumer.ConsumerId)
	if err != nil {
		return nil, e.Wrap("can't add extra segments: ", err)
	}

	for _, segmentId = range extRes {
		_, err = cs.operations.AddOperation(ctx, consumer.ConsumerId, segmentId, actionId)
		if err != nil {
			return nil, respository_errors.ErrCannotCreate
		}
	}

	res = append(res, extRes...)

	return res, nil
}
func (cs *ConsumersService) AddSegmentsToConsumer(ctx context.Context, consumer dto.ConsumerRequest) ([]int, error) {
	var err error
	defer func() {
		err = e.WrapIfErr("Service consumer: ", err)
	}()

	check, err := cs.consumers.GetSegmentsById(ctx, consumer.ConsumerId)
	if len(check) == 0 {
		return nil, ErrUserNotFound
	}

	actionId, err := cs.actions.GetIdByAction(ctx, entity.ActionTypeCreate)
	if err != nil {
		return nil, e.Wrap("can't get action: ", err)
	}

	var consSegmentId int
	var segmentId int
	var res []int
	var id int
	for _, value := range consumer.Segments {
		segmentId, err = cs.segments.GetIdBySegment(ctx, value.SegmentName)
		if err != nil {
			return nil, respository_errors.ErrNotFound
		}

		if value.TTL.IsZero() {
			consSegmentId, err = cs.consumersSegments.AddConsumerSegmentTTL(ctx, segmentId, value.TTL)
			if err != nil {
				return nil, err
			}
		} else {
			consSegmentId, err = cs.consumersSegments.AddConsumerSegment(ctx, segmentId)
			if err != nil {
				return nil, err
			}
		}

		id, err = cs.consumers.CreateConsumer(ctx, consSegmentId, consumer.ConsumerId)
		if err != nil {
			return nil, respository_errors.ErrCannotAdd
		}

		res = append(res, id)

		_, err = cs.operations.AddOperation(ctx, consumer.ConsumerId, consSegmentId, actionId)
		if err != nil {
			return nil, respository_errors.ErrCannotCreate
		}

	}
	return res, nil
}
func (cs *ConsumersService) DeleteSegmentsFromConsumer(ctx context.Context, consumer dto.ConsumerRequest) error {
	var err error
	defer func() {
		err = e.WrapIfErr("Service consumer: ", err)
	}()

	for _, segment := range consumer.Segments {
		err = cs.consumers.DeleteSegmentFromConsumer(ctx, consumer.ConsumerId, segment.SegmentName)
		if err != nil {
			return err
		}
	}
	return nil
}
func (cs *ConsumersService) UpdateSegmentsTTL(ctx context.Context, consumer dto.ConsumerRequest) error {
	var err error
	defer func() {
		err = e.WrapIfErr("Service consumer: ", err)
	}()

	for _, segment := range consumer.Segments {
		err = cs.consumersSegments.UpdateSegmentTTL(ctx, consumer.ConsumerId, segment.SegmentName, segment.TTL)
		if err != nil {
			return err
		}
	}
	return nil
}
func (cs *ConsumersService) GetConsumerSegments(ctx context.Context, consumer dto.ConsumerId) (dto.ConsumerResponse, error) {
	var err error
	defer func() {
		err = e.WrapIfErr("Service consumer: ", err)
	}()

	check, err := cs.consumers.GetSegmentsById(ctx, consumer.ConsumerId)
	if len(check) == 0 {
		return dto.ConsumerResponse{}, ErrUserNotFound
	}

	err = cs.consumersSegments.DeleteExpiredTTL(ctx, consumer.ConsumerId)
	if err != nil {
		return dto.ConsumerResponse{}, e.Wrap("can't delete expired ttl: ", err)
	}

	var result dto.ConsumerResponse
	segments, err := cs.consumers.GetAllSegmentsByConsumerId(ctx, consumer.ConsumerId)
	if err != nil {
		return dto.ConsumerResponse{}, e.Wrap("can't get all segments: ", err)
	}

	result.ConsumerId = consumer.ConsumerId
	for index := range segments {
		result.SegmentsName[index] = segments[index].SegmentName
	}

	return result, nil
}

func (cs *ConsumersService) automaticAdd(ctx context.Context, consumerId int) ([]int, error) {
	var err error
	defer func() {
		err = e.WrapIfErr("Service consumer: ", err)
	}()

	var count int
	count, err = cs.consumers.GetCountConsumers(ctx)
	if err != nil {
		return nil, e.Wrap("can't get count consumers: ", err)
	}
	segments, err := cs.segments.GetAllSegments(ctx)
	if err != nil {
		return nil, e.Wrap("can't get all segments: ", err)
	}

	var needAdd dto.ConsumerRequest
	for index, value := range segments {
		if count_percent.CheckCount(count+1, value.Percent) {
			needAdd.Segments[index].SegmentName = value.Name
		}
	}

	needAdd.ConsumerId = consumerId
	res, err := cs.AddSegmentsToConsumer(ctx, needAdd)
	if err != nil {
		return nil, e.Wrap("can't automatic add segments: ", err)
	}

	return res, nil

}
