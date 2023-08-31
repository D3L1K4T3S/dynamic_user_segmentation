package service

import (
	"context"
	"dynamic-user-segmentation/internal/entity"
	"dynamic-user-segmentation/internal/repository"
	"dynamic-user-segmentation/internal/repository/respository_errors"
	"dynamic-user-segmentation/internal/service/dto"
	"dynamic-user-segmentation/pkg/util/count_percent"
	e "dynamic-user-segmentation/pkg/util/errors"
	"log"
)

type ConsumersService struct {
	consumers         repository.Consumers
	segments          repository.Segments
	consumersSegments repository.ConsumersSegments
	operations        repository.Operations
	actions           repository.Actions
}

func NewConsumersService(consumers repository.Consumers, segments repository.Segments, consumersSegments repository.ConsumersSegments, operations repository.Operations, actions repository.Actions) *ConsumersService {
	return &ConsumersService{
		consumers:         consumers,
		segments:          segments,
		consumersSegments: consumersSegments,
		operations:        operations,
		actions:           actions,
	}
}

func (cs *ConsumersService) CreateConsumer(ctx context.Context, consumer dto.ConsumerRequest) ([]int, error) {
	var err error
	var consSegmentId int
	var segmentId int
	var res []int
	var id int
	defer func() {
		err = e.WrapIfErr("Service consumer: ", err)
	}()

	check, err := cs.consumers.GetSegmentsById(ctx, consumer.ConsumerId)
	if len(check) > 0 {
		return nil, ErrUserAlreadyExists
	}

	var actionId int
	actionId, err = cs.actions.GetIdByAction(ctx, entity.ActionTypeCreate)
	if err != nil {
		return nil, e.Wrap("can't get action: ", err)
	}

	if consumer.Segments == nil {
		id, err = cs.consumers.AddNullSegmentByConsumerId(ctx, consumer.ConsumerId)
		if err != nil {
			return nil, err
		}

		segmentId, err = cs.segments.GetIdBySegment(ctx, dto.Null)
		if err != nil {
			return nil, respository_errors.ErrNotFound
		}

		_, err = cs.operations.AddOperation(ctx, consumer.ConsumerId, segmentId, actionId)
		if err != nil {
			return nil, respository_errors.ErrCannotCreate
		}

		res = append(res, id)
		return res, nil
	}

	for _, value := range consumer.Segments {
		segmentId, err = cs.segments.GetIdBySegment(ctx, value.SegmentName)
		if err != nil {
			return nil, respository_errors.ErrNotFound
		}

		if !value.TTL.IsZero() {
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

		id, err = cs.consumers.CreateConsumer(ctx, consumer.ConsumerId, consSegmentId)
		if err != nil {
			return nil, respository_errors.ErrCannotCreate
		}

		res = append(res, id)

		_, err = cs.operations.AddOperation(ctx, consumer.ConsumerId, segmentId, actionId)
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

	if len(check) == 1 && check[0] == 0 {
		err = cs.consumers.DeleteNullSegmentByConsumerId(ctx, consumer.ConsumerId)
		if err != nil {
			return nil, err
		}
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

		if !value.TTL.IsZero() {
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

		id, err = cs.consumers.CreateConsumer(ctx, consumer.ConsumerId, consSegmentId)
		if err != nil {
			return nil, respository_errors.ErrCannotAdd
		}

		res = append(res, id)

		_, err = cs.operations.AddOperation(ctx, consumer.ConsumerId, segmentId, actionId)
		if err != nil {
			return nil, respository_errors.ErrCannotCreate
		}

	}
	return res, nil
}
func (cs *ConsumersService) DeleteSegmentsFromConsumer(ctx context.Context, consumer dto.ConsumerRequestDelete) error {
	var err error
	defer func() {
		err = e.WrapIfErr("Service consumer: ", err)
	}()

	var segmentId int
	actionId, err := cs.actions.GetIdByAction(ctx, entity.ActionTypeDelete)
	if err != nil {
		return e.Wrap("can't get action: ", err)
	}
	for _, segment := range consumer.Segments {
		err = cs.consumersSegments.DeleteConsumerSegment(ctx, consumer.ConsumerId, segment.SegmentName)
		if err != nil {
			return err
		}

		segmentId, err = cs.segments.GetIdBySegment(ctx, segment.SegmentName)
		if err != nil {
			return ErrSegmentNotFound
		}

		_, err = cs.operations.AddOperation(ctx, consumer.ConsumerId, segmentId, actionId)
		if err != nil {
			return e.Wrap("can't add operation", err)
		}
	}

	check, err := cs.consumers.GetSegmentsById(ctx, consumer.ConsumerId)
	if len(check) == 0 {
		_, err = cs.consumers.AddNullSegmentByConsumerId(ctx, consumer.ConsumerId)
		if err != nil {
			return e.Wrap("can't add null segments", err)
		}
	}

	return nil
}
func (cs *ConsumersService) UpdateSegmentsTTL(ctx context.Context, consumer dto.ConsumerRequest) error {
	var err error
	defer func() {
		err = e.WrapIfErr("Service consumer: ", err)
	}()

	var segmentId int
	actionId, err := cs.actions.GetIdByAction(ctx, entity.ActionTypeUpdate)
	if err != nil {
		return e.Wrap("can't get action: ", err)
	}
	for _, segment := range consumer.Segments {
		err = cs.consumersSegments.UpdateSegmentTTL(ctx, consumer.ConsumerId, segment.SegmentName, segment.TTL)
		if err != nil {
			return err
		}
		segmentId, err = cs.segments.GetIdBySegment(ctx, segment.SegmentName)
		if err != nil {
			return respository_errors.ErrNotFound
		}

		_, err = cs.operations.AddOperation(ctx, consumer.ConsumerId, segmentId, actionId)
		if err != nil {
			return respository_errors.ErrCannotCreate
		}
	}
	return nil
}
func (cs *ConsumersService) GetConsumerSegments(ctx context.Context, consumer dto.ConsumerId) (dto.ConsumerResponse, error) {
	var err error
	defer func() {
		err = e.WrapIfErr("Service consumer: ", err)
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()

	check, err := cs.consumers.GetSegmentsById(ctx, consumer.ConsumerId)
	if len(check) == 0 {
		return dto.ConsumerResponse{}, ErrUserNotFound
	}

	err = cs.consumersSegments.DeleteExpiredTTL(ctx, consumer.ConsumerId)
	if err != nil {
		return dto.ConsumerResponse{}, e.Wrap("can't delete expired ttl: ", err)
	}

	segments, err := cs.consumers.GetAllSegmentsByConsumerId(ctx, consumer.ConsumerId)
	if err != nil {
		return dto.ConsumerResponse{}, e.Wrap("can't get all segments: ", err)
	}

	var result dto.ConsumerResponse
	result.ConsumerId = consumer.ConsumerId
	for index := range segments {
		result.SegmentsName = append(result.SegmentsName, segments[index].SegmentName)
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
