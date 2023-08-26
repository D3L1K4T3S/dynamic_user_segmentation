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

func (or *OperationsRepository) GetOperationsInTime(ctx context.Context, start, end time.Time) ([]entity.Operations, error) {
	var err error
	defer func() {
		err = e.WrapIfErr("problem with get operation in time: ", err)
	}()

	query := "SELECT * FROM operations WHERE created_at >= $1 and created_at <= $2"

	rows, err := or.Pool.Query(ctx, query, start, end)
	if err != nil {
		return nil, e.Wrap("can't query: ", err)
	}

	defer rows.Close()

	var operations []entity.Operations
	for rows.Next() {
		var operation entity.Operations
		err = rows.Scan(&operation.Id, &operation.UserId,
			&operation.SegmentId, &operation.ActionId, &operation.Created)
		if err != nil {
			return nil, e.Wrap("can't scan in struct data: ", err)
		}
		operations = append(operations, operation)
	}

	return operations, nil
}
