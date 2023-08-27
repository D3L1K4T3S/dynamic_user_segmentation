package service

import (
	"context"
	"dynamic-user-segmentation/internal/entity"
	"dynamic-user-segmentation/internal/repository"
)

type ActionsService struct {
	actions repository.Actions
}

func NewActionsService(actionsRepository repository.Actions) *ActionsService {
	return &ActionsService{actions: actionsRepository}
}

func (as *ActionsService) GetActionById(ctx context.Context, id int) (entity.Action, error) {
	return as.actions.GetActionById(ctx, id)
}
func (as *ActionsService) GetIdByAction(ctx context.Context, action entity.Action) (int, error) {
	return as.actions.GetIdByAction(ctx, action)
}
