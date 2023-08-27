package postgresql

import (
	"context"
	"dynamic-user-segmentation/internal/entity"
	"dynamic-user-segmentation/pkg/client/db/postgresql"
)

type ConsumersRepository struct {
	*postgresql.PostgreSQL
}

func NewConsumersRepository(pg *postgresql.PostgreSQL) *ConsumersRepository {
	return &ConsumersRepository{pg}
}

func (cr *ConsumersRepository) CreateConsumer(ctx context.Context, id int, segments ...entity.Segments) (int, error) {

}
func (cr *ConsumersRepository) DeleteConsumer(ctx context.Context, id int) error {

}
func (cr *ConsumersRepository) AddSegmentToConsumer(ctx context.Context, id int, segments ...entity.Segments) (bool, error) {

}
func (cr *ConsumersRepository) DeleteSegmentFromConsumer(ctx context.Context, id int, segments ...entity.Segments) (bool, error) {
	
}
