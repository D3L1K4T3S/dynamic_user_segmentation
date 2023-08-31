package service

import (
	"context"
	"dynamic-user-segmentation/internal/repository"
	"dynamic-user-segmentation/internal/service/dto"
	"dynamic-user-segmentation/pkg/hash"
	"time"
)

type Actions interface {
	CreateAction(ctx context.Context, action dto.ActionsRequest) (int, error)
	DeleteAction(ctx context.Context, action dto.ActionsRequest) error
}

type Auth interface {
	CreateUser(ctx context.Context, au dto.AuthUser) (int, error)
	DeleteUser(ctx context.Context, au dto.AuthUser) error
	GenerateToken(ctx context.Context, user dto.AuthUser) (string, error)
	ParseToken(token string) (int, error)
}

type Consumers interface {
	CreateConsumer(ctx context.Context, consumer dto.ConsumerRequest) ([]int, error)
	AddSegmentsToConsumer(ctx context.Context, consumer dto.ConsumerRequest) ([]int, error)
	DeleteSegmentsFromConsumer(ctx context.Context, consumer dto.ConsumerRequestDelete) error
	UpdateSegmentsTTL(ctx context.Context, consumer dto.ConsumerRequest) error
	GetConsumerSegments(ctx context.Context, consumer dto.ConsumerId) (dto.ConsumerResponse, error)
}

type Operations interface {
	GetHistoryForPeriod(ctx context.Context, request dto.OperationsRequest) (dto.OperationsResponse, error)
	GetReportFile(ctx context.Context, request dto.OperationsRequest) ([]byte, error)
}

type Segments interface {
	CreateSegment(ctx context.Context, segment dto.SegmentsRequest) (int, error)
	DeleteSegment(ctx context.Context, segment dto.SegmentsRequest) error
	UpdateSegment(ctx context.Context, segment dto.SegmentsRequest) error
	GetAllSegments(ctx context.Context) ([]dto.SegmentsResponse, error)
}

type Services struct {
	Actions
	Auth
	Consumers
	Operations
	Segments
}

type ServicesDependencies struct {
	Repository *repository.Repositories
	Hash       hash.PasswordHash

	TokenTTL time.Duration
	SignKey  string
}

func NewServices(servicesDependencies ServicesDependencies) *Services {
	return &Services{
		Actions: NewActionsService(
			servicesDependencies.Repository.Actions),
		Auth: NewAuthService(
			servicesDependencies.Repository.Users,
			servicesDependencies.Hash,
			servicesDependencies.SignKey,
			servicesDependencies.TokenTTL),
		Consumers: NewConsumersService(
			servicesDependencies.Repository.Consumers,
			servicesDependencies.Repository.Segments,
			servicesDependencies.Repository.ConsumersSegments,
			servicesDependencies.Repository.Operations,
			servicesDependencies.Repository.Actions),
		Operations: NewOperationsService(
			servicesDependencies.Repository.Operations,
			servicesDependencies.Repository.ConsumersSegments,
			servicesDependencies.Repository.Segments,
			servicesDependencies.Repository.Actions),
		Segments: NewSegmentsService(
			servicesDependencies.Repository.Segments),
	}
}
