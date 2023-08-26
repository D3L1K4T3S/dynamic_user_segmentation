package postgresql

import (
	"context"
	"dynamic-user-segmentation/internal/entity"
	"dynamic-user-segmentation/pkg/client/db/postgresql"
	e "dynamic-user-segmentation/pkg/util/errors"
	"errors"
	"github.com/jackc/pgx/v5"
)

type ActionsRepository struct {
	*postgresql.PostgreSQL
}

func NewActionsRepository(pg *postgresql.PostgreSQL) *ActionsRepository {
	return &ActionsRepository{pg}
}

func (ar *ActionsRepository) GetActionById(ctx context.Context, id int) (entity.Action, error) {
	var err error
	defer func() {
		err = e.WrapIfErr("problem with get action by id: ", err)
	}()

	query := "SELECT name FROM actions WHERE id = $1"

	var action entity.Action
	err = ar.Pool.QueryRow(ctx, query, id).Scan(&action)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", e.Wrap("can't find action by id in table: ", err)
		}
		return "", err
	}

	return action, nil
}
func (ar *ActionsRepository) GetIdByAction(ctx context.Context, action entity.Action) (int, error) {
	var err error
	defer func() {
		err = e.WrapIfErr("problem in GetIdByAction: ", err)
	}()

	query := "SELECT id FROM actions WHERE name = $1"

	var id int
	err = ar.Pool.QueryRow(ctx, query, action).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, e.Wrap("can't find id by action in table: ", err)
		}
		return 0, err
	}

	return id, nil
}
