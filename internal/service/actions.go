package service

import (
	"context"
	"dynamic-user-segmentation/internal/repository"
	"dynamic-user-segmentation/internal/service/dto"
)

type ActionsService struct {
	actions repository.Actions
}

func NewActionsService(actionsRepository repository.Actions) *ActionsService {
	return &ActionsService{actions: actionsRepository}
}

func (as *ActionsService) CreateAction(ctx context.Context, action dto.ActionsRequest) (int, error) {
	return as.actions.CreateAction(ctx, action.Name)
}
func (as *ActionsService) DeleteAction(ctx context.Context, action dto.ActionsRequest) error {
	return as.actions.DeleteAction(ctx, action.Name)
}

//func (as *ActionsService) GetActionById(ctx context.Context, id int) (dto.ActionsResponse, error) {
//	var actions dto.ActionsResponse
//	res, err := as.actions.GetActionById(ctx, id)
//	if err != nil {
//		return dto.ActionsResponse{}, e.Wrap("can't do service: ", err)
//	}
//	actions.Id = res.Id
//	actions.Name = res.Name
//
//	return actions, nil
//}
//func (as *ActionsService) GetIdByAction(ctx context.Context, action dto.ActionsRequest) (int, error) {
//	return as.actions.GetIdByAction(ctx, action.Name)
//}
