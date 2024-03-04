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

func TestRefreshHandler_Refresh(t *testing.T) {
	method := "PUT"
	URL := "/refresh"

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockUserReauthorizer := mocks.NewMockUserReauthorizer(mockCtrl)
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
			mock: mockUserReauthorizer.EXPECT().
				ReauthUser(dto.ReAuthUserDTO{RefreshToken: "refresh_token"}).
				Return(&dto.ReAuthUserDTO{
					AccessToken:  "new_access_token",
					RefreshToken: "new_refresh_token",
				}, nil).
				Times(1),
			body:     `{"refresh_token":"refresh_token"}`,
			wantCode: http.StatusOK,
			wantBody: "{\"status\":\"OK\",\"access_token\":\"new_access_token\",\"refresh_token\":\"new_refresh_token\"}\n",
		},
		{
			name:     "wrong body",
			mock:     nil,
			body:     `"wrong_tag":"wrong_value"`,
			wantCode: http.StatusBadRequest,
			wantBody: "{\"status\":\"Error\",\"error\":\"failed to decode request\"}\n",
		},
		{
			name: "invalid credentials",
			mock: mockUserReauthorizer.EXPECT().
				ReauthUser(dto.ReAuthUserDTO{RefreshToken: "refresh_token"}).
				Return(nil, service.ErrInvalidCredentials).
				Times(1),
			body:     `{"refresh_token":"refresh_token"}`,
			wantCode: http.StatusUnauthorized,
			wantBody: "{\"status\":\"Error\",\"error\":\"invalid credentials\"}\n",
		},
		{
			name: "failed to authorize",
			mock: mockUserReauthorizer.EXPECT().
				ReauthUser(dto.ReAuthUserDTO{RefreshToken: "refresh_token"}).
				Return(nil, errors.New("any")).
				Times(1),
			body:     `{"refresh_token":"refresh_token"}`,
			wantCode: http.StatusInternalServerError,
			wantBody: "{\"status\":\"Error\",\"error\":\"failed to update refresh token\"}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(method, URL, bytes.NewBufferString(tt.body))

			handler := auth.Refresh(log, mockUserReauthorizer)
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
