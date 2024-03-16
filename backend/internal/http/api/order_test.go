package api_test

import (
	"cloud-render/internal/dto"
	"cloud-render/internal/http/api"
	mocks "cloud-render/internal/mocks/service/api"
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

func TestOrderHandler_Order(t *testing.T) {
	method := "GET"
	URL := "/orders/{id}"

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockOneOrderProvider := mocks.NewMockOneOrderProivder(mockCtrl)
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
			mock: mockOneOrderProvider.EXPECT().
				GetOneOrder(int64(1)).
				Return(&dto.GetOrderDTO{
					OrderStatus:  "status",
					DownloadLink: "link",
				}, nil).
				Times(1),
			id:       "1",
			wantCode: http.StatusOK,
			wantBody: "{\"status\":\"OK\",\"order_status\":\"status\",\"download_link\":\"link\"}\n",
		},
		{
			name: "correct no link",
			mock: mockOneOrderProvider.EXPECT().
				GetOneOrder(int64(1)).
				Return(&dto.GetOrderDTO{
					OrderStatus:  "status",
					DownloadLink: "",
				}, nil).
				Times(1),
			id:       "1",
			wantCode: http.StatusOK,
			wantBody: "{\"status\":\"OK\",\"order_status\":\"status\"}\n",
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
			wantBody: "{\"status\":\"Error\",\"error\":\"invalid id param\"}\n",
		},
		{
			name: "order not found",
			mock: mockOneOrderProvider.EXPECT().
				GetOneOrder(int64(1)).
				Return(nil, service.ErrOrderNotFound).
				Times(1),
			id:       "1",
			wantCode: http.StatusNotFound,
			wantBody: "{\"status\":\"Error\",\"error\":\"order not found\"}\n",
		},
		{
			name: "failed to get order",
			mock: mockOneOrderProvider.EXPECT().
				GetOneOrder(int64(1)).
				Return(nil, errors.New("unknown")).
				Times(1),
			id:       "1",
			wantCode: http.StatusInternalServerError,
			wantBody: "{\"status\":\"Error\",\"error\":\"failed to get order\"}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(method, URL, nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.id)

			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			handler := api.Order(log, mockOneOrderProvider)
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
