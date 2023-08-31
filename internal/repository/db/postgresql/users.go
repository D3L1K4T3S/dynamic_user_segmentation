package postgresql

import (
	"context"
	"dynamic-user-segmentation/internal/entity"
	"dynamic-user-segmentation/internal/repository/respository_errors"
	"dynamic-user-segmentation/pkg/client/db/postgresql"
	e "dynamic-user-segmentation/pkg/util/errors"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type UsersRepository struct {
	*postgresql.PostgreSQL
}

func NewUsersRepository(pg *postgresql.PostgreSQL) *UsersRepository {
	return &UsersRepository{pg}
}

func (ur *UsersRepository) CreateUser(ctx context.Context, user entity.Users) (int, error) {
	var err error
	defer func() {
		err = e.WrapIfErr(respository_errors.RepositoryPostgresMsg, err)
	}()

	query := "INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id"

	var id int
	err = ur.Pool.QueryRow(ctx, query, user.Username, user.Password).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == "23505" {
				return 0, respository_errors.ErrAlreadyExists
			}
		}
		return 0, e.Wrap(respository_errors.CannotDoQueryMsg, err)
	}

	return id, nil
}
func (ur *UsersRepository) DeleteUser(ctx context.Context, user entity.Users) error {
	var err error
	defer func() {
		err = e.WrapIfErr(respository_errors.RepositoryPostgresMsg, err)
	}()

	query := "DELETE FROM users WHERE username = $1 and password = $2"

	_, err = ur.Pool.Exec(ctx, query, user.Username, user.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return respository_errors.ErrNotFound
		}
		return e.Wrap(respository_errors.CannotDoQueryMsg, err)
	}
	return nil

}
func (ur *UsersRepository) GetUserByUsername(ctx context.Context, username string) (entity.Users, error) {
	var err error
	defer func() {
		err = e.WrapIfErr(respository_errors.RepositoryPostgresMsg, err)
	}()

	query := "SELECT id, password FROM users WHERE username = $1"

	var user entity.Users
	user.Username = username
	err = ur.Pool.QueryRow(ctx, query, username).Scan(&user.Id, &user.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Users{}, respository_errors.ErrNotFound
		}
		return entity.Users{}, e.Wrap(respository_errors.CannotDoQueryMsg, err)
	}

	return user, nil
}
