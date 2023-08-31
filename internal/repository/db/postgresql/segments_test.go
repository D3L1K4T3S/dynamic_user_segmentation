package postgresql

import (
	"context"
	"dynamic-user-segmentation/pkg/client/db/postgresql"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSegmentsRepository_CreateSegment(t *testing.T) {
	type args struct {
		ctx     context.Context
		segment string
		percent float64
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
				ctx:     context.Background(),
				segment: "test_segment",
				percent: 10.0,
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				rows := pgxmock.NewRows([]string{"id"}).AddRow(1)
				m.ExpectQuery("INSERT INTO segments").WithArgs(args.segment, args.percent).WillReturnRows(rows)
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "Segment already exists",
			args: args{
				ctx:     context.Background(),
				segment: "test_segment",
				percent: 10.0,
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectQuery("INSERT INTO segments").
					WithArgs(args.segment, args.percent).
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
				ctx:     context.Background(),
				segment: "test_segment",
				percent: 10.0,
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectQuery("INSERT INTO segments").
					WithArgs(args.segment, args.percent).
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

			segmentsRepositoryMock := NewSegmentsRepository(postgresMock)

			got, err := segmentsRepositoryMock.CreateSegment(tc.args.ctx, tc.args.segment, tc.args.percent)
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

func TestSegmentsRepository_DeleteSegment(t *testing.T) {
	type args struct {
		ctx context.Context
		id  int
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
			name: "Segment already exists",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectQuery("DELETE FROM segments").
					WithArgs(args.id).
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
				id:  1,
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectQuery("DELETE FROM segments").
					WithArgs(args.id).
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

			segmentsRepositoryMock := NewSegmentsRepository(postgresMock)

			err := segmentsRepositoryMock.DeleteSegment(tc.args.ctx, tc.args.id)
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

func TestSegmentsRepository_UpdateSegment(t *testing.T) {
	type args struct {
		ctx     context.Context
		id      int
		percent float64
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
			name: "Segment already exists",
			args: args{
				ctx: context.Background(),
				id:  1,
				percent: 10
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectQuery("UPDATE segments SET percent WHERE id").
					WithArgs(args.percent, args.id).
			},
			want:    1,
			wantErr: true,
		},
		{
			name: "Unexpected error",
			args: args{
				ctx: context.Background(),
				id:  1,
			},
			mockBehavior: func(m pgxmock.PgxPoolIface, args args) {
				m.ExpectQuery("DELETE FROM segments").
					WithArgs(args.id).
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

			segmentsRepositoryMock := NewSegmentsRepository(postgresMock)

			err := segmentsRepositoryMock.UpdateSegment(tc.args.ctx, tc.args.id, tc.args.percent)
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

func TestSegmentsRepository_GetIdBySegment(t *testing.T) {

}

func TestSegmentsRepository_GetAllSegments(t *testing.T) {

}
