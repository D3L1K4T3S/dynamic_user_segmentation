package postgresql

import (
	"context"
	"dynamic-user-segmentation/internal/repository"
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

func (sr *SegmentsRepository) CreateSegment(ctx context.Context, segment string, percent int) (int, error) {
	var err error
	defer func() {
		err = e.WrapIfErr("can't create a new segment", err)
	}()

	query := "INSERT INTO segments (name, percent) VALUES ($1, $2) RETURNING id"

	var id int
	err = sr.Pool.QueryRow(ctx, query, segment, percent).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == "23505" {
				return 0, repository.ErrAlreadyExists
			}
		}
		return 0, e.Wrap("can't do a query: ", err)
	}

	return id, nil
}
func (sr *SegmentsRepository) DeleteSegment(ctx context.Context, id int) error {
	var err error
	defer func() {
		err = e.WrapIfErr("can't delete from segments: ", err)
	}()

	query := "DELETE FROM segments WHERE id = $1"
	_, err = sr.Pool.Exec(ctx, query, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repository.ErrNotFound
		}
		return e.Wrap("can't do a query: ", err)
	}
	return nil
}
func (sr *SegmentsRepository) UpdateSegment(ctx context.Context, id int, percent int) error {
	var err error
	defer func() {
		err = e.WrapIfErr("can't update segment: ", err)
	}()

	query := "UPDATE segments SET percent = $1 WHERE id = $2"
	_, err = sr.Pool.Exec(ctx, query, percent, id)
	if err != nil {
		return e.Wrap("can't do a query: ", err)
	}
	return nil
}
func (sr *SegmentsRepository) GetSegmentById(ctx context.Context, id int) (string, error) {
	var err error
	defer func() {
		err = e.WrapIfErr("can't get segment by id: ", err)
	}()

	query := "SELECT name FROM segments WHERE id = $1"

	var name string
	err = sr.Pool.QueryRow(ctx, query, id).Scan(&name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", repository.ErrNotFound
		}
		return "", e.Wrap("can't do query: ", err)
	}

	return name, nil
}
func (sr *SegmentsRepository) GetIdBySegment(ctx context.Context, segment string) (int, error) {
	var err error
	defer func() {
		err = e.WrapIfErr("can't get id by name segment: ", err)
	}()

	query := "SELECT name FROM segments WHERE id = $1"

	var id int
	err = sr.Pool.QueryRow(ctx, query, segment).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, repository.ErrNotFound
		}
		return 0, e.Wrap("can't do a query: ", err)
	}

	return id, nil
}
