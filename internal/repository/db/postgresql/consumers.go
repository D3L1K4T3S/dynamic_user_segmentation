package postgresql

import (
	"context"
	"dynamic-user-segmentation/internal/entity"
	"dynamic-user-segmentation/internal/repository/respository_errors"
	"dynamic-user-segmentation/pkg/client/db/postgresql"
	e "dynamic-user-segmentation/pkg/util/errors"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
)

type ConsumersRepository struct {
	*postgresql.PostgreSQL
}

func NewConsumersRepository(pg *postgresql.PostgreSQL) *ConsumersRepository {
	return &ConsumersRepository{pg}
}

func (cr *ConsumersRepository) CreateConsumer(ctx context.Context, consumerId int, segmentId int) (int, error) {
	var err error
	defer func() {
		err = e.WrapIfErr("problem with create a consumer: ", err)
	}()

	query := "INSERT INTO consumers(consumer_id, segment_id) VALUES ($1, $2) RETURNING id"

	var id int
	err = cr.Pool.QueryRow(ctx, query, consumerId, segmentId).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == "23505" {
				return 0, respository_errors.ErrAlreadyExists
			}
		}
		return 0, e.Wrap("can't do a query: ", err)
	}

	return id, nil
}

func (cr *ConsumersRepository) GetSegmentsById(ctx context.Context, consumerId int) ([]int, error) {
	var err error
	defer func() {
		err = e.WrapIfErr("problem with get segments consumer: ", err)
	}()

	query := "SELECT segment_id FROM consumers WHERE consumer_id = $1"

	rows, err := cr.Pool.Query(ctx, query, consumerId)
	if err != nil {
		return nil, e.Wrap("can't do a query: ", err)
	}

	defer rows.Close()

	segmentsId := make([]int, 0)
	for rows.Next() {
		var segmentId int
		err = rows.Scan(&segmentId)
		if err != nil {
			return nil, e.Wrap("can't scan a values from rows: ", err)
		}
		segmentsId = append(segmentsId, segmentId)
	}
	return segmentsId, nil
}

func (cr *ConsumersRepository) DeleteNullSegmentByConsumerId(ctx context.Context, consumerId int) error {
	var err error
	defer func() {
		err = e.WrapIfErr("Repository postgres: ", err)
	}()

	query := "DELETE FROM consumers WHERE segment_id IS NULL and consumer_id = $1"
	_, err = cr.Pool.Exec(ctx, query, consumerId)
	if err != nil {
		return e.Wrap(respository_errors.CannotDoQueryMsg, err)
	}
	return nil
}

func (cr *ConsumersRepository) GetCountConsumers(ctx context.Context) (int, error) {
	var err error
	defer func() {
		err = e.WrapIfErr("Repository postgres: ", err)
	}()

	query := "SELECT COUNT(DISTINCT consumer_id) FROM consumers"

	var count int
	err = cr.Pool.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, e.Wrap("can't do a query: ", err)
	}
	return count, nil
}

func (cr *ConsumersRepository) GetAllSegmentsByConsumerId(ctx context.Context, consumerId int) ([]entity.ComplexConsumerSegments, error) {
	var err error
	defer func() {
		err = e.WrapIfErr("Postgres client: ", err)
	}()

	query := "SELECT consumer_id, segments.name FROM consumers LEFT JOIN consumers_segments ON consumers.segment_id = consumers_segments.id LEFT JOIN segments ON consumers_segments.segment_id = segments.id WHERE consumers.consumer_id = $1"

	rows, err := cr.Pool.Query(ctx, query, consumerId)
	if err != nil {
		return nil, e.Wrap("can't do a query: ", err)
	}
	defer rows.Close()

	var result []entity.ComplexConsumerSegments
	for rows.Next() {
		var cs entity.ComplexConsumerSegments
		err = rows.Scan(&cs.ConsumerId, &cs.SegmentName)
		if err != nil {
			return nil, e.Wrap("can't scan rows: ", err)
		}
		result = append(result, cs)
	}
	return result, nil
}

func (cr *ConsumersRepository) AddNullSegmentByConsumerId(ctx context.Context, consumerId int) (int, error) {
	var err error
	defer func() {
		err = e.WrapIfErr(respository_errors.RepositoryPostgresMsg, err)
	}()

	query := "INSERT INTO consumers(consumer_id) VALUES ($1) RETURNING id"

	var id int
	err = cr.Pool.QueryRow(ctx, query, consumerId).Scan(&id)
	if err != nil {
		return 0, e.Wrap(respository_errors.CannotDoQueryMsg, err)
	}
	return id, nil
}

func (cr *ConsumersRepository) ExistConsumer(ctx context.Context, consumerId int) (bool, error) {
	var err error
	defer func() {
		err = e.WrapIfErr(respository_errors.RepositoryPostgresMsg, err)
	}()

	query := "SELECT EXISTS (SELECT * FROM consumers WHERE consumers.consumer_id = $1)"

	var res bool
	err = cr.Pool.QueryRow(ctx, query, consumerId).Scan(&res)
	if err != nil {
		return false, e.Wrap(respository_errors.CannotDoQueryMsg, err)
	}
	return res, nil

}
