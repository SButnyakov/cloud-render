package api_test

import (
	"cloud-render/internal/http/api"
	"cloud-render/internal/lib/config"
	mocks "cloud-render/internal/mocks/service/api"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
)

func TestSubscribeHandler_Orders(t *testing.T) {
	method := "POST"
	URL := "/subscribe"

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockUserSubscriber := mocks.NewMockUserSubscriber(mockCtrl)
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	tests := []struct {
		name     string
		mock     *gomock.Call
		id       string
		wantCode int
		wantBody string
	}{
		{
			name: "correct",
			mock: mockUserSubscriber.EXPECT().
				SubscribeUser(int64(1)).
				Return(nil).
				Times(1),
			id:       "1",
			wantCode: http.StatusOK,
			wantBody: "{\"status\":\"OK\"}\n",
		},
		{
			name: "failed to subscribe",
			mock: mockUserSubscriber.EXPECT().
				SubscribeUser(int64(1)).
				Return(errors.New("any")).
				Times(1),
			id:       "1",
			wantCode: http.StatusInternalServerError,
			wantBody: "{\"status\":\"Error\",\"error\":\"failed to subscribe\"}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(method, URL, nil)

			ctx := context.WithValue(r.Context(), "uid", int64(1))
			r = r.WithContext(ctx)

			handler := api.Subscribe(log, &config.Config{}, mockUserSubscriber)
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
