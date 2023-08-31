package postgresql

import (
	"context"
	"dynamic-user-segmentation/internal/entity"
	"dynamic-user-segmentation/internal/repository/respository_errors"
	"dynamic-user-segmentation/pkg/client/db/postgresql"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUsersRepository_CreateUser(t *testing.T) {
	type args struct {
		ctx  context.Context
		user entity.Users
	}

	type MockBehavior func(pgm pgxmock.PgxPoolIface, args args)

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
				user: entity.Users{
					Username: "test_user",
					Password: "test_password",
				},
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				rows := pgxmock.NewRows([]string{"id"}).AddRow(1)
				m.ExpectQuery("INSERT INTO users").WithArgs(args.user.Username, args.user.Password).WillReturnRows(rows)
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "User already exists",
			args: args{
				ctx: context.Background(),
				user: entity.Users{
					Username: "test_user",
					Password: "test_password",
				},
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectQuery("INSERT INTO users").
					WithArgs(args.user.Username, args.user.Password).
					WillReturnError(&pgconn.PgError{
						Code: "23505",
					})
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "Unexpected error",
			args: args{
				ctx: context.Background(),
				user: entity.Users{
					Username: "test_user",
					Password: "test_password",
				},
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectQuery("INSERT INTO users").
					WithArgs(args.user.Username, args.user.Password).
					WillReturnError(errors.New("unexpected error"))
			},
			want:    0,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			poolMock, _ := pgxmock.NewPool()
			defer poolMock.Close()
			tc.mockBehavior(poolMock, tc.args)

			postgresMock := &postgresql.PostgreSQL{
				Pool: poolMock,
			}

			usersRepositoryMock := NewUsersRepository(postgresMock)

			got, err := usersRepositoryMock.CreateUser(tc.args.ctx, tc.args.user)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.want, got)

			err = poolMock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func TestUsersRepository_GetUserByUsernameAndPassword(t *testing.T) {
	type args struct {
		ctx      context.Context
		username string
	}

	type MockBehavior func(pgm pgxmock.PgxPoolIface, args args)

	testCases := []struct {
		name         string
		args         args
		mockBehavior MockBehavior
		want         entity.Users
		wantErr      bool
	}{
		{
			name: "OK",
			args: args{
				ctx:      context.Background(),
				username: "test_user",
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				rows := pgxmock.NewRows([]string{"id", "password"}).
					AddRow(1, "test_password")

				m.ExpectQuery("SELECT id, password FROM users").
					WithArgs(args.username).
					WillReturnRows(rows)
			},
			want: entity.Users{
				Id:       1,
				Username: "test_user",
				Password: "test_password",
			},
			wantErr: false,
		},
		{
			name: "User not found",
			args: args{
				ctx:      context.Background(),
				username: "unknown_user",
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {

				m.ExpectQuery("SELECT id, password FROM users").
					WithArgs(args.username).
					WillReturnError(pgx.ErrNoRows)
			},
			want:    entity.Users{},
			wantErr: true,
		},
		{
			name: "Unexpected error",
			args: args{
				ctx:      context.Background(),
				username: "unknown_name",
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {

				m.ExpectQuery("SELECT id, password FROM users").
					WithArgs(args.username).WillReturnError(errors.New("unexpected error"))
			},
			want:    entity.Users{},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			poolMock, _ := pgxmock.NewPool()
			defer poolMock.Close()
			tc.mockBehavior(poolMock, tc.args)

			postgresMock := &postgresql.PostgreSQL{
				Pool: poolMock,
			}

			usersRepositoryMock := NewUsersRepository(postgresMock)

			got, err := usersRepositoryMock.GetUserByUsername(tc.args.ctx, tc.args.username)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.want, got)

			err = poolMock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func TestUsersRepository_DeleteUser(t *testing.T) {
	type args struct {
		ctx  context.Context
		user entity.Users
	}

	type MockBehavior func(pgm pgxmock.PgxPoolIface, args args)

	testCases := []struct {
		name         string
		args         args
		mockBehavior MockBehavior
		want         error
		wantErr      bool
	}{
		{
			name: "User not found",
			args: args{
				ctx: context.Background(),
				user: entity.Users{
					Username: "unknown_user",
					Password: "unknown_password",
				},
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectQuery("DELETE FROM users").
					WithArgs(args.user.Username, args.user.Password).
					WillReturnError(respository_errors.ErrNotFound)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Unexpected error",
			args: args{
				ctx: context.Background(),
				user: entity.Users{
					Username: "unknown_user",
					Password: "unknown_password",
				},
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectQuery("DELETE FROM users").
					WithArgs(args.user.Username, args.user.Password).
					WillReturnError(errors.New("unexpected error"))
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			poolMock, _ := pgxmock.NewPool()
			defer poolMock.Close()
			tc.mockBehavior(poolMock, tc.args)

			postgresMock := &postgresql.PostgreSQL{
				Pool: poolMock,
			}

			usersRepositoryMock := NewUsersRepository(postgresMock)

			err := usersRepositoryMock.DeleteUser(tc.args.ctx, tc.args.user)
			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.want, err)

			err = poolMock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}
