package api_test

import (
	"cloud-render/internal/dto"
	"cloud-render/internal/http/api"
	mocks "cloud-render/internal/mocks/service/api"
	"cloud-render/internal/service"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

func TestUserHandler_User(t *testing.T) {
	method := "GET"
	URL := "/user"

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockUserInfoProvider := mocks.NewMockUserInfoProvider(mockCtrl)
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	now := time.Now()
	today := now.Format("02-01-2006")

	tests := []struct {
		name     string
		mock     *gomock.Call
		id       string
		wantCode int
		wantBody string
	}{
		{
			name: "correct with exp date",
			mock: mockUserInfoProvider.EXPECT().
				GetExpireDateWithUserInfo(int64(1)).
				Return(&dto.UserInfoDTO{
					Login:          "login",
					Email:          "email@gmail.com",
					ExpirationDate: &now,
				}, nil).
				Times(1),
			id:       "1",
			wantCode: http.StatusOK,
			wantBody: fmt.Sprintf("{\"status\":\"OK\",\"login\":\"login\",\"email\":\"email@gmail.com\",\"expirationDate\":\"%s\"}\n", today),
		},
		{
			name: "correct no exp date",
			mock: mockUserInfoProvider.EXPECT().
				GetExpireDateWithUserInfo(int64(1)).
				Return(&dto.UserInfoDTO{
					Login:          "login",
					Email:          "email@gmail.com",
					ExpirationDate: nil,
				}, nil).
				Times(1),
			id:       "1",
			wantCode: http.StatusOK,
			wantBody: "{\"status\":\"OK\",\"login\":\"login\",\"email\":\"email@gmail.com\",\"expirationDate\":null}\n",
		},
		{
			name: "failed to fetch user info",
			mock: mockUserInfoProvider.EXPECT().
				GetExpireDateWithUserInfo(int64(1)).
				Return(nil, service.ErrExternalError).
				Times(1),
			id:       "1",
			wantCode: http.StatusBadRequest,
			wantBody: "{\"status\":\"Error\",\"error\":\"failed to fetch user info\"}\n",
		},
		{
			name: "failed to get subscription info",
			mock: mockUserInfoProvider.EXPECT().
				GetExpireDateWithUserInfo(int64(1)).
				Return(nil, service.ErrFailedToGetSubscription).
				Times(1),
			id:       "1",
			wantCode: http.StatusInternalServerError,
			wantBody: "{\"status\":\"Error\",\"error\":\"failed to get subcription info\"}\n",
		},
		{
			name: "failed to get subscription info",
			mock: mockUserInfoProvider.EXPECT().
				GetExpireDateWithUserInfo(int64(1)).
				Return(nil, errors.New("any")).
				Times(1),
			id:       "1",
			wantCode: http.StatusInternalServerError,
			wantBody: "{\"status\":\"Error\",\"error\":\"failed to get user data\"}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(method, URL, nil)

			ctx := context.WithValue(r.Context(), "uid", int64(1))
			r = r.WithContext(ctx)

			handler := api.User(log, mockUserInfoProvider)
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
