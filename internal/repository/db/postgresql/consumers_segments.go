package postgresql

import (
	"context"
	"dynamic-user-segmentation/internal/repository"
	"dynamic-user-segmentation/pkg/client/db/postgresql"
	e "dynamic-user-segmentation/pkg/util/errors"
	"errors"
	"github.com/jackc/pgx/v5"
	"time"
)

type ConsumersSegmentsRepository struct {
	*postgresql.PostgreSQL
}

func NewConsumersSegmentsRepository(pg *postgresql.PostgreSQL) *ConsumersSegmentsRepository {
	return &ConsumersSegmentsRepository{pg}
}

func (csr *ConsumersSegmentsRepository) AddConsumerSegment(ctx context.Context, segmentId int, TTL time.Time) (int, error) {
	var err error
	defer func() {
		err = e.WrapIfErr("can't add segment to consumer", err)
	}()

	query := "INSERT INTO consumers_segments (segment_id, ttl) VALUES ($1,$2) RETURNING id"

	var id int
	err = csr.Pool.QueryRow(ctx, query, segmentId, TTL).Scan(&id)
	if err != nil {
		return 0, e.Wrap("can't do a query", err)
	}
	return id, nil
}
func (csr *ConsumersSegmentsRepository) DeleteConsumerSegment(ctx context.Context, segmentId int) error {
	var err error
	defer func() {
		err = e.WrapIfErr("can't delete segment from consumer", err)
	}()

	query := "DELETE FROM consumers_segments WHERE id = $1"

	_, err = csr.Pool.Exec(ctx, query, segmentId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repository.ErrNotFound
		}
		return e.Wrap("can't do a query: ", err)
	}

	return nil
}
func (csr *ConsumersSegmentsRepository) UpdateSegmentTTL(ctx context.Context, segmentId int, TTL time.Time) error {
	var err error
	defer func() {
		err = e.WrapIfErr("can't update ttl in segment: ", err)
	}()

	query := "UPDATE consumers_segments SET ttl = $1 WHERE id = $2"
	_, err = csr.Pool.Exec(ctx, query, TTL, segmentId)
	if err != nil {
		return e.Wrap("can't do a query: ", err)
	}
	return nil
}
