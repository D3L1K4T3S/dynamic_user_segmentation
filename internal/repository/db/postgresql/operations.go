package postgresql

import (
	"context"
	"dynamic-user-segmentation/internal/entity"
	"dynamic-user-segmentation/pkg/client/db/postgresql"
	e "dynamic-user-segmentation/pkg/util/errors"
	"time"
)

type OperationsRepository struct {
	*postgresql.PostgreSQL
}

func NewOperationsRepository(pg *postgresql.PostgreSQL) *OperationsRepository {
	return &OperationsRepository{pg}
}

func (or *OperationsRepository) GetOperationsInTime(ctx context.Context, consumerId int, start, end time.Time) ([]entity.ComplexOperations, error) {
	var err error
	defer func() {
		err = e.WrapIfErr("problem with get operation in time: ", err)
	}()

	query := "SELECT consumer_id, segments.name, actions.name, created_at FROM OPERATIONS LEFT JOIN SEGMENTS ON OPERATIONS.segment_id = segments.id LEFT JOIN ACTIONS ON OPERATIONS.action_id = ACTIONS.id WHERE OPERATIONS.consumer_id = $1 AND created_at >= $2 AND created_at =< $3;"

	rows, err := or.Pool.Query(ctx, query, consumerId, start, end)
	if err != nil {
		return nil, e.Wrap("can't query: ", err)
	}

	defer rows.Close()

	var operations []entity.ComplexOperations
	for rows.Next() {
		var operation entity.ComplexOperations
		err = rows.Scan(&operation.ConsumerId, &operation.SegmentName, &operation.ActionName, &operation.Created)
		if err != nil {
			return nil, e.Wrap("can't scan in struct data: ", err)
		}
		operations = append(operations, operation)
	}

	return operations, nil
}
func (or *OperationsRepository) AddOperation(ctx context.Context, consumerId int, segmentId int, actionId int) (int, error) {
	var err error
	defer func() {
		err = e.WrapIfErr("Postgresql operations: ", err)
	}()

	query := "INSERT INTO operations VALUES ($1,$2,$3,$4) RETURNING id"

	var id int
	err = or.Pool.QueryRow(ctx, query, consumerId, segmentId, actionId).Scan(&id)
	if err != nil {
		return 0, e.Wrap("can't do a query: ", err)
	}

	return id, nil
}
