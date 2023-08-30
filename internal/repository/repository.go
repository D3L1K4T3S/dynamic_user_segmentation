package repository

import (
	"context"
	"dynamic-user-segmentation/internal/entity"
	postgres "dynamic-user-segmentation/internal/repository/db/postgresql"
	"dynamic-user-segmentation/pkg/client/db/postgresql"
	"time"
)

type Actions interface {
	CreateAction(ctx context.Context, action string) (int, error)
	DeleteAction(ctx context.Context, action string) error
	//GetActionById(ctx context.Context, actionId int) (string, error)
	GetIdByAction(ctx context.Context, action string) (int, error)
}

type Consumers interface {
	CreateConsumer(ctx context.Context, consumerId int, segmentId int) (int, error)
	DeleteConsumer(ctx context.Context, consumerId int) error
	AddSegmentToConsumer(ctx context.Context, id int, segmentId int) error
	GetSegmentsById(ctx context.Context, consumerId int) ([]int, error)
	DeleteSegmentFromConsumer(ctx context.Context, consumerId int, segmentName string) error
	GetCountConsumers(ctx context.Context) (int, error)
	GetAllSegmentsByConsumerId(ctx context.Context, consumerId int) ([]entity.ComplexConsumerSegments, error)
}

type ConsumersSegments interface {
	AddConsumerSegmentTTL(ctx context.Context, segmentId int, TTL time.Time) (int, error)
	AddConsumerSegment(ctx context.Context, segmentId int) (int, error)
	DeleteConsumerSegment(ctx context.Context, consumerId int, segmentName string) error
	GetSegmentIdById(ctx context.Context, id int) (int, error)
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
	GetSegmentById(ctx context.Context, id int) (string, error)
	GetIdBySegment(ctx context.Context, segment string) (int, error)
	GetAllSegments(ctx context.Context) ([]entity.Segments, error)
}

type Users interface {
	CreateUser(ctx context.Context, user entity.Users) (int, error)
	DeleteUser(ctx context.Context, user entity.Users) error
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
