package repository

import (
	"context"
	"dynamic-user-segmentation/internal/entity"
	postgres "dynamic-user-segmentation/internal/repository/db/postgresql"
	"dynamic-user-segmentation/pkg/client/db/postgresql"
	"time"
)

type Actions interface {
	CreateAction(ctx context.Context, action entity.Action) (int, error)
	DeleteAction(ctx context.Context, action string) error
	GetActionById(ctx context.Context, actionId int) (entity.Action, error)
	GetIdByAction(ctx context.Context, action entity.Action) (int, error)
}

type Consumers interface {
	CreateConsumer(ctx context.Context, consumerId int) (int, error)
	DeleteConsumer(ctx context.Context, consumerId int) error
	AddSegmentToConsumer(ctx context.Context, id int, segmentId int) error
	GetSegmentsById(ctx context.Context, consumerId int) ([]int, error)
	DeleteSegmentFromConsumer(ctx context.Context, id int, segmentId int) error
}

type ConsumersSegments interface {
	AddConsumerSegment(ctx context.Context, segmentId int, TTL time.Time) (int, error)
	DeleteConsumerSegment(ctx context.Context, segmentId int) error
	UpdateSegmentTTL(ctx context.Context, segmentId int, TTL time.Time) error
}

type Operations interface {
	GetOperationsInTime(ctx context.Context, start, end time.Time) ([]entity.Operations, error)
}

type Segments interface {
	CreateSegment(ctx context.Context, segment string, percent int) (int, error)
	DeleteSegment(ctx context.Context, id int) error
	UpdateSegment(ctx context.Context, id int, percent int) error
	GetSegmentById(ctx context.Context, id int) (string, error)
	GetIdBySegment(ctx context.Context, segment string) (int, error)
}

type Users interface {
	CreateUser(ctx context.Context, username string, password string) (int, error)
	GetUserByID(ctx context.Context, id int) (entity.Users, error)
	GetUserByUsername(ctx context.Context, username string) (entity.Users, error)
	GetIdByUsername(ctx context.Context, username string) (int, error)
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
		Actions:           postgres.NewActionsRepository(pg),
		Consumers:         postgres.NewConsumersRepository(pg),
		ConsumersSegments: postgres.NewConsumersSegmentsRepository(pg),
		Operations:        postgres.NewOperationsRepository(pg),
		Segments:          postgres.NewSegmentsRepository(pg),
		Users:             postgres.NewUsersRepository(pg),
	}
}
