package auth_test

import (
	"cloud-render/internal/dto"
	"cloud-render/internal/http/auth"
	mocks "cloud-render/internal/mocks/service/auth"
	"cloud-render/internal/service"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
)

func TestInfoHandler_Info(t *testing.T) {
	method := "GET"
	URL := "/info/{id}"

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockUserProvider := mocks.NewMockUserProvider(mockCtrl)
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	tests := []struct {
		name     string
		mock     *gomock.Call
		id       string
		wantCode int
		wantBody string
	}{
		{
			name: "correct with link",
			mock: mockUserProvider.EXPECT().
				GetUser(int64(1)).
				Return(&dto.GetUserDTO{
					Login: "login",
					Email: "email@gmail.com",
				}, nil).
				Times(1),
			id:       "1",
			wantCode: http.StatusOK,
			wantBody: "{\"status\":\"OK\",\"login\":\"login\",\"email\":\"email@gmail.com\"}\n",
		},
		{
			name:     "no param provided",
			mock:     nil,
			id:       "",
			wantCode: http.StatusBadRequest,
			wantBody: "{\"status\":\"Error\",\"error\":\"empty id param\"}\n",
		},
		{
			name:     "invalid id param",
			mock:     nil,
			id:       "not int",
			wantCode: http.StatusBadRequest,
			wantBody: "{\"status\":\"Error\",\"error\":\"invalid id\"}\n",
		},
		{
			name: "order not found",
			mock: mockUserProvider.EXPECT().
				GetUser(int64(1)).
				Return(nil, service.ErrUserNotFound).
				Times(1),
			id:       "1",
			wantCode: http.StatusNotFound,
			wantBody: "{\"status\":\"Error\",\"error\":\"user not found\"}\n",
		},
		{
			name: "failed to get order",
			mock: mockUserProvider.EXPECT().
				GetUser(int64(1)).
				Return(nil, errors.New("any")).
				Times(1),
			id:       "1",
			wantCode: http.StatusInternalServerError,
			wantBody: "{\"status\":\"Error\",\"error\":\"failed to get user\"}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(method, URL, nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.id)

			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			handler := auth.Info(log, mockUserProvider)
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
