package postgresql

import (
	"context"
	"dynamic-user-segmentation/internal/repository/respository_errors"
	"dynamic-user-segmentation/pkg/client/db/postgresql"
	e "dynamic-user-segmentation/pkg/util/errors"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type ActionsRepository struct {
	*postgresql.PostgreSQL
}

func NewActionsRepository(pg *postgresql.PostgreSQL) *ActionsRepository {
	return &ActionsRepository{pg}
}

func (ar *ActionsRepository) CreateAction(ctx context.Context, action string) (int, error) {
	var err error
	defer func() {
		err = e.WrapIfErr("problem with add action: ", err)
	}()

	query := "INSERT INTO actions (name) VALUES ($1) RETURNING id"

	var id int
	err = ar.Pool.QueryRow(ctx, query, action).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == "23505" {
				return 0, respository_errors.ErrAlreadyExists
			}
		}
		return 0, e.Wrap("can't do a query: ", err)
	}

	return id, nil
}
func (ar *ActionsRepository) DeleteAction(ctx context.Context, action string) error {
	var err error
	defer func() {
		err = e.WrapIfErr("problem in delete action", err)
	}()

	query := "DELETE FROM actions WHERE name = $1"
	_, err = ar.Pool.Exec(ctx, query, action)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return respository_errors.ErrNotFound
		}
		return e.Wrap("can't do a query: ", err)
	}
	return nil
}
func (ar *ActionsRepository) GetIdByAction(ctx context.Context, action string) (int, error) {
	var err error
	defer func() {
		err = e.WrapIfErr("problem in GetIdByAction: ", err)
	}()

	query := "SELECT id FROM actions WHERE name = $1"

	var id int
	err = ar.Pool.QueryRow(ctx, query, action).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, respository_errors.ErrNotFound
		}
		return 0, err
	}

	return id, nil
}

//func (ar *ActionsRepository) GetActionById(ctx context.Context, actionId int) (string, error) {
//	var err error
//	defer func() {
//		err = e.WrapIfErr("problem with get action by id: ", err)
//	}()
//
//	query := "SELECT name FROM actions WHERE id = $1"
//
//	var action string
//	err = ar.Pool.QueryRow(ctx, query, actionId).Scan(&action)
//	if err != nil {
//		if errors.Is(err, pgx.ErrNoRows) {
//			return "", repository.ErrNotFound
//		}
//		return "", err
//	}
//	return action, nil
//}
