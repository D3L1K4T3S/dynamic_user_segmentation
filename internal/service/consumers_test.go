package service

import (
	"context"
	rmocks "dynamic-user-segmentation/internal/mocks/repository"
	"dynamic-user-segmentation/internal/service/dto"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConsumersService_AddSegmentsToConsumer(t *testing.T) {
	type args struct {
		ctx   context.Context
		input dto.ActionsRequest
	}

	type MockBehavior func(a *rmocks.MockActions, args args)

	testCases := []struct {
		name         string
		args         args
		mockBehavior MockBehavior
		want         int
		wantErr      bool
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				input: dto.ActionsRequest{
					Name: "create",
				},
			},
			mockBehavior: func(a *rmocks.MockActions, args args) {
				a.EXPECT().CreateAction(args.ctx, args.input.Name).Return(
					1, nil)
			},
			want:    1,
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			controller := gomock.NewController(t)
			defer controller.Finish()

			operationRepo := rmocks.NewMockActions(controller)
			tc.mockBehavior(operationRepo, tc.args)

			service := NewActionsService(operationRepo)

			got, err := service.CreateAction(tc.args.ctx, tc.args.input)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestConsumersService_CreateConsumer(t *testing.T) {
	type args struct {
		ctx   context.Context
		input dto.ConsumerRequest
	}

	type MockBehavior func(a *rmocks.MockActions, args args)

	testCases := []struct {
		name         string
		args         args
		mockBehavior MockBehavior
		want         int
		wantErr      bool
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				input: dto.ConsumerRequest{
					ConsumerId: 1,
					Segments: []dto.DataSegments{{"discount"}, "message", ""},
				}
			},
			mockBehavior: func(a *rmocks.MockActions, args args) {
				a.EXPECT().CreateAction(args.ctx, args.input.Name).Return(
					1, nil)
			},
			want:    1,
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			controller := gomock.NewController(t)
			defer controller.Finish()

			operationRepo := rmocks.NewMockActions(controller)
			tc.mockBehavior(operationRepo, tc.args)

			service := NewActionsService(operationRepo)

			got, err := service.CreateAction(tc.args.ctx, tc.args.input)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestConsumersService_UpdateSegmentsTTL(t *testing.T) {

}


func TestConsumersService_GetConsumerSegments(t *testing.T) {

}