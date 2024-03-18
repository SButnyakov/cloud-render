package buffer_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"cloud-render/internal/dto"
	"cloud-render/internal/http/buffer"
	mocks "cloud-render/internal/mocks/service/buffer"
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/golang/mock/gomock"
)

func TestUpdateHandler(t *testing.T) {
	method := "PUT"
	URL := "/{uid}/blend/update/{filename}/{status}"

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockOrderStatusUpdater := mocks.NewMockOrderStatusUpdater(mockCtrl)
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	tests := []struct {
		name     string
		id       string
		filename string
		status   string
		wantCode int
		wantBody string
		mock     *gomock.Call
	}{
		{
			name:     "Successful update",
			id:       "1",
			filename: "test_file.txt",
			status:   "status",
			wantCode: http.StatusOK,
			wantBody: "{\"status\":\"OK\"}\n",
			mock: mockOrderStatusUpdater.EXPECT().
				UpdateOrderStatus(dto.UpdateOrderStatusDTO{
					OrderId: int64(1),
					Status:  "status",
				}).
				Return(nil).
				Times(1),
		},
		{
			name:     "Invalid user ID",
			id:       "not int",
			filename: "test_file.txt",
			status:   "status",
			wantCode: http.StatusBadRequest,
			wantBody: "{\"status\":\"Error\",\"error\":\"invalid user id\"}\n",
		},
		{
			name:     "Error from updater",
			id:       "1",
			filename: "test_file.txt",
			status:   "status",
			wantCode: http.StatusInternalServerError,
			wantBody: "{\"status\":\"Error\",\"error\":\"failed to update status\"}\n",
			mock: mockOrderStatusUpdater.EXPECT().
				UpdateOrderStatus(dto.UpdateOrderStatusDTO{
					OrderId: int64(1),
					Status:  "status",
				}).
				Return(errors.New("any")).
				Times(1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(method, URL, nil)

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("uid", tt.id)
			rctx.URLParams.Add("filename", tt.filename)
			rctx.URLParams.Add("status", tt.status)

			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			handler := buffer.Update(log, mockOrderStatusUpdater)
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
