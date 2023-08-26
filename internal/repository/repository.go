package repository

import (
	"context"
	"dynamic-user-segmentation/internal/entity"
	postgres "dynamic-user-segmentation/internal/repository/db/postgresql"
	"dynamic-user-segmentation/pkg/client/db/postgresql"
	"time"
)

type Actions interface {
	GetActionById(ctx context.Context, id int) (entity.Action, error)
	GetIdByAction(ctx context.Context, action entity.Action) (int, error)
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
	CreateUser(ctx context.Context, userId int) (int, error)
	AddSegmentToUser(ctx context.Context, userId int, segments ...entity.Segments) (bool, error)
	DeleteSegmentFromUser(ctx context.Context, userId int, segments ...entity.Segments) (bool, error)
}

type UsersSegments interface {
	UpdateSegmentTTL(ctx context.Context, segmentID int, TTL time.Time) error
}

type Repositories struct {
	Actions
	Operations
	Segments
	Users
	UsersSegments
}

func NewRepositories(pg *postgresql.PostgreSQL) *Repositories {
	return &Repositories{
		Actions:       postgres.NewActionsRepository(pg),
		Operations:    postgres.NewOperationsRepository(pg),
		Segments:      postgres.NewSegmentsRepository(pg),
		Users:         postgres.NewUsersRepository(pg),
		UsersSegments: postgres.NewUsersSegmentsRepository(pg),
	}
}
