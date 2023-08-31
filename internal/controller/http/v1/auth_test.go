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

func TestAuthRoutes_SignIn(t *testing.T) {
	type args struct {
		ctx   context.Context
		input dto.AuthUser
	}

	type MockBehavior func(m *smocks.MockAuth, args args)

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
				input: dto.AuthUser{
					Username: "test_name",
					Password: "test_password",
				},
			},
			inputBody: `{"username":"test_name","password":"test_password"}`,
			mockBehavior: func(m *smocks.MockAuth, args args) {
				m.EXPECT().GenerateToken(args.ctx, args.input).Return("token", nil)
			},
			wantStatusCode:  200,
			wantRequestBody: `{"token":"token"}` + "\n",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			auth := smocks.NewMockAuth(controller)
			tc.mockBehavior(auth, tc.args)
			services := &service.Services{Auth: auth}

			e := echo.New()
			g := e.Group("/auth")
			newAuthRoutes(g, services.Auth)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/auth/sign-in", bytes.NewBufferString(tc.inputBody))
			r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			e.ServeHTTP(w, r)

			assert.Equal(t, tc.wantStatusCode, w.Code)
			assert.Equal(t, tc.wantRequestBody, w.Body.String())
		})
	}
}
