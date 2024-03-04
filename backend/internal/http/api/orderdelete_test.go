package api_test

import (
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

func TestOrderDeleteHandler_DeleteOrder(t *testing.T) {
	method := "POST"
	URL := "/orders/{id}/delete"

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockDeleteOrderProvider := mocks.NewMockOneOrderSoftDelter(mockCtrl)
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
			mock: mockDeleteOrderProvider.EXPECT().
				SoftDeleteOneOrder(int64(1)).
				Return(nil).
				Times(1),
			id:       "1",
			wantCode: http.StatusOK,
			wantBody: "{\"status\":\"OK\"}\n",
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
			mock: mockDeleteOrderProvider.EXPECT().
				SoftDeleteOneOrder(int64(1)).
				Return(service.ErrOrderNotFound).
				Times(1),
			id:       "1",
			wantCode: http.StatusNotFound,
			wantBody: "{\"status\":\"Error\",\"error\":\"order not found\"}\n",
		},
		{
			name: "failed to delete order",
			mock: mockDeleteOrderProvider.EXPECT().
				SoftDeleteOneOrder(int64(1)).
				Return(errors.New("unknown")).
				Times(1),
			id:       "1",
			wantCode: http.StatusInternalServerError,
			wantBody: "{\"status\":\"Error\",\"error\":\"failed to delete order\"}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(method, URL, nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.id)

			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			handler := api.DeleteOrder(log, mockDeleteOrderProvider)
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
