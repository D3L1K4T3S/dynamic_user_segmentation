package service

import (
	"context"
	"dynamic-user-segmentation/internal/repository"
	"dynamic-user-segmentation/internal/service/dbo"
	e "dynamic-user-segmentation/pkg/util/errors"
)

type ActionsService struct {
	actions repository.Actions
}

func NewActionsService(actionsRepository repository.Actions) *ActionsService {
	return &ActionsService{actions: actionsRepository}
}

func (as *ActionsService) CreateAction(ctx context.Context, action dbo.ActionsRequest) (int, error) {
	return as.actions.CreateAction(ctx, action.Name)
}
func (as *ActionsService) DeleteAction(ctx context.Context, action dbo.ActionsRequest) error {
	return as.actions.DeleteAction(ctx, action.Name)
}
func (as *ActionsService) GetActionById(ctx context.Context, id int) (dbo.ActionsResponse, error) {
	var actions dbo.ActionsResponse
	res, err := as.actions.GetActionById(ctx, id)
	if err != nil {
		return dbo.ActionsResponse{}, e.Wrap("can't do service: ", err)
	}
	actions.Id = res.Id
	actions.Name = res.Name

	return actions, nil
}
func (as *ActionsService) GetIdByAction(ctx context.Context, action dbo.ActionsRequest) (int, error) {
	return as.actions.GetIdByAction(ctx, action.Name)
}
