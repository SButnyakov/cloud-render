package auth_test

import (
	"bytes"
	"cloud-render/internal/dto"
	"cloud-render/internal/http/auth"
	mocks "cloud-render/internal/mocks/service/auth"
	"cloud-render/internal/service"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestSignInHandler_SignIn(t *testing.T) {
	method := "POST"
	URL := "/signin"

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockUserAuthorizer := mocks.NewMockUserAuthorizer(mockCtrl)
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	tests := []struct {
		name     string
		mock     *gomock.Call
		body     string
		wantCode int
		wantBody string
	}{
		{
			name: "correct",
			mock: mockUserAuthorizer.EXPECT().
				AuthUser(dto.AuthUserDTO{
					LoginOrEmail: "login",
					Password:     "password",
				}).
				Return(&dto.AuthUserDTO{
					AccessToken:  "new_access_token",
					RefreshToken: "new_refresh_token",
				}, nil).
				Times(1),
			body:     `{"login_or_email":"login", "password":"password"}`,
			wantCode: http.StatusOK,
			wantBody: "{\"status\":\"OK\",\"access_token\":\"new_access_token\",\"refresh_token\":\"new_refresh_token\"}\n",
		},
		{
			name:     "empty body",
			mock:     nil,
			body:     ``,
			wantCode: http.StatusBadRequest,
			wantBody: "{\"status\":\"Error\",\"error\":\"empty request\"}\n",
		},
		{
			name:     "no login",
			mock:     nil,
			body:     `{"password":"password"}`,
			wantCode: http.StatusBadRequest,
			wantBody: "{\"status\":\"Error\",\"error\":\"field LoginOrEmail is a required field\"}\n",
		},
		{
			name:     "no password",
			mock:     nil,
			body:     `{"login_or_email":"login"}`,
			wantCode: http.StatusBadRequest,
			wantBody: "{\"status\":\"Error\",\"error\":\"field Password is a required field\"}\n",
		},
		{
			name:     "empty password",
			mock:     nil,
			body:     `{"login_or_email":"login", "password":""}`,
			wantCode: http.StatusBadRequest,
			wantBody: "{\"status\":\"Error\",\"error\":\"field Password is a required field\"}\n",
		},
		{
			name:     "empty login",
			mock:     nil,
			body:     `{"login_or_email":"", "password":"password"}`,
			wantCode: http.StatusBadRequest,
			wantBody: "{\"status\":\"Error\",\"error\":\"field LoginOrEmail is a required field\"}\n",
		},
		{
			name: "invalid credentials",
			mock: mockUserAuthorizer.EXPECT().
				AuthUser(dto.AuthUserDTO{
					LoginOrEmail: "login",
					Password:     "password",
				}).
				Return(nil, service.ErrInvalidCredentials).
				Times(1),
			body:     `{"login_or_email":"login", "password":"password"}`,
			wantCode: http.StatusBadRequest,
			wantBody: "{\"status\":\"Error\",\"error\":\"invalid credentials\"}\n",
		},
		{
			name: "server-side authorization failed",
			mock: mockUserAuthorizer.EXPECT().
				AuthUser(dto.AuthUserDTO{
					LoginOrEmail: "login",
					Password:     "password",
				}).
				Return(nil, errors.New("any")).
				Times(1),
			body:     `{"login_or_email":"login", "password":"password"}`,
			wantCode: http.StatusInternalServerError,
			wantBody: "{\"status\":\"Error\",\"error\":\"server-side authorization failed\"}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(method, URL, bytes.NewBufferString(tt.body))

			handler := auth.SignIn(log, mockUserAuthorizer)
			handler(w, r)

			if w.Result().StatusCode != tt.wantCode {
				t.Errorf("status code should be [%d] but received [%d]", tt.wantCode, w.Result().StatusCode)
			}

			if w.Body.String() != tt.wantBody {
				t.Errorf("the response body should be [%s] but received [%s]", tt.wantBody, w.Body.String())
			}
		})
	}
}
