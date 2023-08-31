package v1

import (
	"bytes"
	"context"
	smocks "dynamic-user-segmentation/internal/mocks/service"
	"dynamic-user-segmentation/internal/service"
	"dynamic-user-segmentation/internal/service/dto"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestActionsRepository_CreateAction(t *testing.T) {
	type args struct {
		ctx   context.Context
		input dto.ActionsRequest
	}

	type MockBehavior func(m *smocks.MockActions, args args)

	testCases := []struct {
		name            string
		args            args
		inputBody       string
		mockBehavior    MockBehavior
		wantStatusCode  int
		wantRequestBody string
	}{
		{
			name: "OK",
			args: args{
				ctx: context.Background(),
				input: dto.ActionsRequest{
					Name: "create",
				},
			},
			inputBody: `{"name":"create"}`,
			mockBehavior: func(m *smocks.MockActions, args args) {
				m.EXPECT().CreateAction(args.ctx, args.input).Return(1, nil)
			},
			wantStatusCode:  201,
			wantRequestBody: `{"id":1}` + "\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			actions := smocks.NewMockActions(controller)
			tc.mockBehavior(actions, tc.args)
			services := &service.Services{Actions: actions}

			e := echo.New()
			g := e.Group("/api/v1/")

			newActionsRoutes(g, services.Actions)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/api/v1/actions/create", bytes.NewBufferString(tc.inputBody))
			r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			e.ServeHTTP(w, r)

			assert.Equal(t, tc.wantStatusCode, w.Code)
			assert.Equal(t, tc.wantRequestBody, w.Body.String())
		})
	}
}
