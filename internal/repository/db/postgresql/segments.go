package postgresql

import (
	"context"
	"dynamic-user-segmentation/internal/entity"
	"dynamic-user-segmentation/internal/repository/respository_errors"
	"dynamic-user-segmentation/pkg/client/db/postgresql"
	e "dynamic-user-segmentation/pkg/util/errors"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type SegmentsRepository struct {
	*postgresql.PostgreSQL
}

func NewSegmentsRepository(pg *postgresql.PostgreSQL) *SegmentsRepository {
	return &SegmentsRepository{pg}
}

func (sr *SegmentsRepository) CreateSegment(ctx context.Context, segment string, percent float64) (int, error) {
	var err error
	defer func() {
		err = e.WrapIfErr(respository_errors.RepositoryPostgresMsg, err)
	}()

	query := "INSERT INTO segments (name, percent) VALUES ($1, $2) RETURNING id"

	var id int
	err = sr.Pool.QueryRow(ctx, query, segment, percent).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == "23505" {
				return 0, respository_errors.ErrAlreadyExists
			}
		}
		return 0, e.Wrap(respository_errors.CannotDoQueryMsg, err)
	}

	return id, nil
}
func (sr *SegmentsRepository) DeleteSegment(ctx context.Context, id int) error {
	var err error
	defer func() {
		err = e.WrapIfErr(respository_errors.RepositoryPostgresMsg, err)
	}()

	query := "DELETE FROM segments WHERE id = $1"
	_, err = sr.Pool.Exec(ctx, query, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return respository_errors.ErrNotFound
		}
		return e.Wrap(respository_errors.CannotDoQueryMsg, err)
	}
	return nil
}
func (sr *SegmentsRepository) UpdateSegment(ctx context.Context, id int, percent float64) error {
	var err error
	defer func() {
		err = e.WrapIfErr(respository_errors.RepositoryPostgresMsg, err)
	}()

	query := "UPDATE segments SET percent = $1 WHERE id = $2"
	_, err = sr.Pool.Exec(ctx, query, percent, id)
	if err != nil {
		return e.Wrap(respository_errors.CannotDoQueryMsg, err)
	}
	return nil
}
func (sr *SegmentsRepository) GetIdBySegment(ctx context.Context, segment string) (int, error) {
	var err error
	defer func() {
		err = e.WrapIfErr(respository_errors.RepositoryPostgresMsg, err)
	}()

	query := "SELECT id FROM segments WHERE name = $1"

	var id int
	err = sr.Pool.QueryRow(ctx, query, segment).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, respository_errors.ErrNotFound
		}
		return 0, e.Wrap(respository_errors.CannotDoQueryMsg, err)
	}

	return id, nil
}
func (sr *SegmentsRepository) GetAllSegments(ctx context.Context) ([]entity.Segments, error) {
	var err error
	defer func() {
		err = e.WrapIfErr(respository_errors.RepositoryPostgresMsg, err)
	}()

	query := "SELECT * FROM segments"

	rows, err := sr.Pool.Query(ctx, query)
	if err != nil {
		return nil, e.Wrap("can't do a query: ", err)
	}

	defer rows.Close()

	var segments []entity.Segments
	for rows.Next() {
		var segment entity.Segments
		err = rows.Scan(&segment.Id, &segment.Name, &segment.Percent)
		if err != nil {
			return nil, e.Wrap(respository_errors.CannotDoQueryMsg, err)
		}
		segments = append(segments, segment)
	}

	return segments, nil
}

func (sr *SegmentsRepository) ExistSegmentConsumer(ctx context.Context, consumerId int, segmentName string) (bool, error) {
	var err error
	defer func() {
		err = e.WrapIfErr("Repository segments: ", err)
	}()

	query := "SELECT EXISTS (SELECT* FROM consumers LEFT JOIN CONSUMERS_SEGMENTS ON consumers.segment_id = CONSUMERS_SEGMENTS.id LEFT JOIN SEGMENTS ON CONSUMERS_SEGMENTS.segment_id = SEGMENTS.id WHERE SEGMENTS.name = $2 and CONSUMERS.consumer_id = $1)"

	var res bool
	err = sr.Pool.QueryRow(ctx, query, consumerId, segmentName).Scan(&res)
	if err != nil {
		return false, e.Wrap("can't exist segment: ", err)
	}
	return res, nil
}
