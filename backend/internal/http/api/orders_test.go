package api_test

import (
	"cloud-render/internal/dto"
	"cloud-render/internal/http/api"
	mocks "cloud-render/internal/mocks/service/api"
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

func TestOrdersHandler_Orders(t *testing.T) {
	method := "GET"
	URL := "/orders"

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockManyOrdersProvider := mocks.NewMockManyOrdersProvider(mockCtrl)
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
			name: "correct",
			mock: mockManyOrdersProvider.EXPECT().
				GetManyOrders(int64(1)).
				Return([]dto.GetOrderDTO{
					{
						Id:           int64(1),
						Filename:     "filename1",
						Date:         now,
						OrderStatus:  "status1",
						DownloadLink: "link1",
					},
					{
						Id:          int64(2),
						Filename:    "filename2",
						Date:        now,
						OrderStatus: "status2",
					},
				}, nil).
				Times(1),
			id:       "1",
			wantCode: http.StatusOK,
			wantBody: fmt.Sprintf("{\"status\":\"OK\",\"orders\":[{\"id\":1,\"filename\":\"filename1\",\"date\":\"%s\",\"status\":\"status1\",\"downloadLink\":\"link1\"},{\"id\":2,\"filename\":\"filename2\",\"date\":\"%s\",\"status\":\"status2\"}]}\n", today, today),
		},
		{
			name: "correct empty",
			mock: mockManyOrdersProvider.EXPECT().
				GetManyOrders(int64(1)).
				Return([]dto.GetOrderDTO{}, nil).
				Times(1),
			id:       "1",
			wantCode: http.StatusOK,
			wantBody: "{\"status\":\"OK\",\"orders\":[]}\n",
		},
		{
			name: "failed to get orders",
			mock: mockManyOrdersProvider.EXPECT().
				GetManyOrders(int64(1)).
				Return(nil, errors.New("any")).
				Times(1),
			id:       "1",
			wantCode: http.StatusInternalServerError,
			wantBody: "{\"status\":\"Error\",\"error\":\"failed to get orders\"}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(method, URL, nil)

			ctx := context.WithValue(r.Context(), "uid", int64(1))
			r = r.WithContext(ctx)

			handler := api.Orders(log, mockManyOrdersProvider)
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
