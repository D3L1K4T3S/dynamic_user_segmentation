package service

import (
	"bytes"
	"context"
	"dynamic-user-segmentation/internal/repository"
	"dynamic-user-segmentation/internal/service/dto"
	e "dynamic-user-segmentation/pkg/util/errors"
	"encoding/csv"
	"strconv"
)

type OperationsService struct {
	operations       repository.Operations
	consumerSegments repository.ConsumersSegments
	segments         repository.Segments
	actions          repository.Actions
}

func NewOperationsService(operations repository.Operations, consumerSegments repository.ConsumersSegments, segments repository.Segments, actions repository.Actions) *OperationsService {
	return &OperationsService{
		operations:       operations,
		consumerSegments: consumerSegments,
		segments:         segments,
		actions:          actions,
	}
}

func (os *OperationsService) GetReportFile(ctx context.Context, request dto.OperationsRequest) ([]byte, error) {
	operations, err := os.operations.GetOperationsInTime(ctx, request.ConsumerId, request.StartAt, request.EndAt)
	if err != nil {
		return nil, e.Wrap("can't get operations in time: ", err)
	}

	data := bytes.Buffer{}
	writer := csv.NewWriter(&data)

	for index := range operations {
		err = writer.Write([]string{strconv.Itoa(operations[index].ConsumerId),
			operations[index].SegmentName,
			operations[index].ActionName,
			operations[index].Created.String()})
		if err != nil {
			return nil, e.Wrap("can't write csv: ", err)
		}
	}

	writer.Flush()
	if err = writer.Error(); err != nil {
		return nil, e.Wrap("can't write csv: ", err)
	}

	return data.Bytes(), nil
}

func (os *OperationsService) GetHistoryForPeriod(ctx context.Context, request dto.OperationsRequest) (dto.OperationsResponse, error) {
	var err error
	defer func() {
		err = e.WrapIfErr("Service operations: ", err)
	}()

	res, err := os.operations.GetOperationsInTime(ctx, request.ConsumerId, request.StartAt, request.EndAt)
	if err != nil {
		return dto.OperationsResponse{}, e.Wrap("can't get operations: ", err)
	}

	var operations dto.OperationsResponse
	for _, v := range res {
		var operation dto.OperationsData
		operation.ActionName = v.ActionName
		operation.SegmentName = v.SegmentName
		operation.Date = v.Created
		operations.OperationsData = append(operations.OperationsData, operation)
	}
	operations.ConsumerId = request.ConsumerId

	return operations, nil
}
