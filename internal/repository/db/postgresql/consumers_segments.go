package postgresql

import (
	"context"
	"dynamic-user-segmentation/pkg/client/db/postgresql"
	e "dynamic-user-segmentation/pkg/util/errors"
	"time"
)

type ConsumersSegmentsRepository struct {
	*postgresql.PostgreSQL
}

func NewConsumersSegmentsRepository(pg *postgresql.PostgreSQL) *ConsumersSegmentsRepository {
	return &ConsumersSegmentsRepository{pg}
}

func (usr *ConsumersSegmentsRepository) AddConsumerSegment(ctx context.Context, segmentId int, TTL time.Time) (int, error) {

}
func (usr *ConsumersSegmentsRepository) DeleteConsumerSegment(ctx context.Context, segmentId int) error {

}

func (usr *ConsumersSegmentsRepository) UpdateSegmentTTL(ctx context.Context, segmentID int, TTL time.Time) error {
	var err error
	defer func() {
		err = e.WrapIfErr("can't update ttl in segment: ", err)
	}()

	query := "UPDATE users_segments SET ttl = $1 WHERE id = $2"
	_, err = usr.Pool.Exec(ctx, query, TTL, segmentID)
	if err != nil {
		return e.Wrap("can't do a query: ", err)
	}
	return nil
}
