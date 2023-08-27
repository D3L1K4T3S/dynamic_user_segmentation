package postgresql

import (
	"context"
	"dynamic-user-segmentation/internal/entity"
	"dynamic-user-segmentation/internal/repository"
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
		err = e.WrapIfErr("problem in create user: ", err)
	}()

	query := "INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id"

	var id int
	err = ur.Pool.QueryRow(ctx, query, user.Username, user.Password).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == "23505" {
				return 0, repository.ErrAlreadyExists
			}
		}
		return 0, e.Wrap("can't do a query: ", err)
	}

	return id, nil
}
func (ur *UsersRepository) GetUserByID(ctx context.Context, id int) (entity.Users, error) {
	var err error
	defer func() {
		err = e.WrapIfErr("problem in get user by id: ", err)
	}()

	query := "SELECT username, password FROM users WHERE id = $1"

	var user entity.Users
	user.Id = id
	err = ur.Pool.QueryRow(ctx, query, id).Scan(&user.Username, &user.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Users{}, repository.ErrNotFound
		}
		return entity.Users{}, err
	}

	return user, nil
}
func (ur *UsersRepository) GetUserByUsername(ctx context.Context, username string) (entity.Users, error) {
	var err error
	defer func() {
		err = e.WrapIfErr("problem in get user by id: ", err)
	}()

	query := "SELECT id, password FROM users WHERE username = $1"

	var user entity.Users
	user.Username = username
	err = ur.Pool.QueryRow(ctx, query, username).Scan(&user.Id, &user.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Users{}, repository.ErrNotFound
		}
		return entity.Users{}, err
	}

	return user, nil
}
func (ur *UsersRepository) GetIdByUsername(ctx context.Context, username string) (int, error) {
	var err error
	defer func() {
		err = e.WrapIfErr("problem in get user by id: ", err)
	}()

	query := "SELECT id FROM users WHERE username = $1"

	var id int
	err = ur.Pool.QueryRow(ctx, query, username).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, repository.ErrNotFound
		}
		return 0, err
	}

	return id, nil
}
