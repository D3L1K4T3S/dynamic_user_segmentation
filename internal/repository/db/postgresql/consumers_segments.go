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

func (csr *ConsumersSegmentsRepository) AddConsumerSegmentTTL(ctx context.Context, segmentId int, TTL time.Time) (int, error) {
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
func (csr *ConsumersSegmentsRepository) AddConsumerSegment(ctx context.Context, segmentId int) (int, error) {
	var err error
	defer func() {
		err = e.WrapIfErr("can't add segment to consumer", err)
	}()

	query := "INSERT INTO consumers_segments (segment_id) VALUES ($1) RETURNING id"

	var id int
	err = csr.Pool.QueryRow(ctx, query, segmentId).Scan(&id)
	if err != nil {
		return 0, e.Wrap("can't do a query", err)
	}
	return id, nil
}
func (csr *ConsumersSegmentsRepository) DeleteConsumerSegment(ctx context.Context, consumerId int, segmentName string) error {
	var err error
	defer func() {
		err = e.WrapIfErr("can't delete segment from consumer", err)
	}()

	query := "DELETE FROM CONSUMERS_SEGMENTS WHERE CONSUMERS_SEGMENTS.id IN (SELECT CONSUMERS_SEGMENTS.id FROM CONSUMERS LEFT JOIN CONSUMERS_SEGMENTS ON CONSUMERS.segment_id = CONSUMERS_SEGMENTS.id LEFT JOIN SEGMENTS ON CONSUMERS_SEGMENTS.segment_id = SEGMENTS.id WHERE consumer_id = $1 and name = $2);"

	_, err = csr.Pool.Exec(ctx, query, consumerId, segmentName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repository.ErrNotFound
		}
		return e.Wrap("can't do a query: ", err)
	}

	return nil
}
func (csr *ConsumersSegmentsRepository) UpdateSegmentTTL(ctx context.Context, consumerId int, segmentName string, TTL time.Time) error {
	var err error
	defer func() {
		err = e.WrapIfErr("can't update ttl in segment: ", err)
	}()

	query := "UPDATE CONSUMERS_SEGMENTS SET ttl = $1 WHERE CONSUMERS_SEGMENTS.id IN (SELECT CONSUMERS_SEGMENTS.id FROM CONSUMERS LEFT JOIN CONSUMERS_SEGMENTS ON CONSUMERS.segment_id = CONSUMERS_SEGMENTS.id LEFT JOIN SEGMENTS ON CONSUMERS_SEGMENTS.segment_id = SEGMENTS.id WHERE consumer_id = $2 and name = $3);"
	_, err = csr.Pool.Exec(ctx, query, TTL, consumerId, segmentName)
	if err != nil {
		return e.Wrap("can't do a query: ", err)
	}
	return nil
}
func (csr *ConsumersSegmentsRepository) GetSegmentIdById(ctx context.Context, id int) (int, error) {
	var err error
	defer func() {
		err = e.WrapIfErr("Repository consumer segments: ", err)
	}()

	query := "SELECT segment_id FROM consumers_segments WHERE id = $1"

	var segmentId int
	err = csr.Pool.QueryRow(ctx, query, id).Scan(&segmentId)
	if err != nil {
		return 0, repository.ErrNotFound
	}
	return id, nil
}
func (csr *ConsumersSegmentsRepository) DeleteExpiredTTL(ctx context.Context, consumerId int) error {
	var err error
	defer func() {
		err = e.WrapIfErr("Repository consumers segments: ", err)
	}()

	query := "SELECT delete_expired_ttl($1)"
	_, err = csr.Pool.Exec(ctx, query, consumerId)
	if err != nil {
		return e.Wrap("can't do a query: ", err)
	}
	return nil
}
