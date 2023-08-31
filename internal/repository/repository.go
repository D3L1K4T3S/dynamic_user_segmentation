package repository

import (
	"context"
	"dynamic-user-segmentation/internal/entity"
	postgresql2 "dynamic-user-segmentation/internal/repository/db/postgresql"
	"dynamic-user-segmentation/pkg/client/db/postgresql"
	"time"
)

type Actions interface {
	CreateAction(ctx context.Context, action string) (int, error)
	DeleteAction(ctx context.Context, action string) error
	GetIdByAction(ctx context.Context, action string) (int, error)
}

type Consumers interface {
	CreateConsumer(ctx context.Context, consumerId int, segmentId int) (int, error)
	GetSegmentsById(ctx context.Context, consumerId int) ([]int, error)
	GetCountConsumers(ctx context.Context) (int, error)
	GetAllSegmentsByConsumerId(ctx context.Context, consumerId int) ([]entity.ComplexConsumerSegments, error)
	DeleteNullSegmentByConsumerId(ctx context.Context, consumerId int) error
	AddNullSegmentByConsumerId(ctx context.Context, consumerId int) (int, error)
	ExistConsumer(ctx context.Context, consumerId int) (bool, error)
}

type ConsumersSegments interface {
	AddConsumerSegmentTTL(ctx context.Context, segmentId int, TTL time.Time) (int, error)
	AddConsumerSegment(ctx context.Context, segmentId int) (int, error)
	DeleteConsumerSegment(ctx context.Context, consumerId int, segmentName string) error
	UpdateSegmentTTL(ctx context.Context, consumerId int, segmentName string, TTL time.Time) error
	DeleteExpiredTTL(ctx context.Context, consumerId int) error
}

type Operations interface {
	GetOperationsInTime(ctx context.Context, consumerId int, start, end time.Time) ([]entity.ComplexOperations, error)
	AddOperation(ctx context.Context, consumerId int, segmentId int, actionId int) (int, error)
}

type Segments interface {
	CreateSegment(ctx context.Context, segment string, percent float64) (int, error)
	DeleteSegment(ctx context.Context, id int) error
	UpdateSegment(ctx context.Context, id int, percent float64) error
	GetIdBySegment(ctx context.Context, segment string) (int, error)
	GetAllSegments(ctx context.Context) ([]entity.Segments, error)
	ExistSegmentConsumer(ctx context.Context, consumerId int, segmentName string) (bool, error)
}

type Users interface {
	CreateUser(ctx context.Context, user entity.Users) (int, error)
	DeleteUser(ctx context.Context, user entity.Users) error
	GetUserByUsername(ctx context.Context, username string) (entity.Users, error)
}

type Repositories struct {
	Actions
	Consumers
	ConsumersSegments
	Operations
	Segments
	Users
}

func NewRepositories(pg *postgresql.PostgreSQL) *Repositories {
	return &Repositories{
		Actions:           postgresql2.NewActionsRepository(pg),
		Consumers:         postgresql2.NewConsumersRepository(pg),
		ConsumersSegments: postgresql2.NewConsumersSegmentsRepository(pg),
		Operations:        postgresql2.NewOperationsRepository(pg),
		Segments:          postgresql2.NewSegmentsRepository(pg),
		Users:             postgresql2.NewUsersRepository(pg),
	}
}
