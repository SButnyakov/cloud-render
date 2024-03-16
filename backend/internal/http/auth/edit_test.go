package auth_test

import (
	"bytes"
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

	"github.com/golang/mock/gomock"
)

func TestEditHandler_Edit(t *testing.T) {
	method := "POST"
	URL := "/edit"

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockUserEditor := mocks.NewMockUserEditor(mockCtrl)
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	tests := []struct {
		name     string
		mock     *gomock.Call
		body     string
		uid      int64
		wantCode int
		wantBody string
	}{
		{
			name: "correct",
			mock: mockUserEditor.EXPECT().
				EditUser(dto.EditUserDTO{
					Id:       1,
					Login:    "newlogin",
					Email:    "test@example.com",
					Password: "newpassword",
				}).
				Return(nil).
				Times(1),
			body:     `{"login":"newlogin", "email":"test@example.com", "password":"newpassword"}`,
			uid:      1,
			wantCode: http.StatusOK,
			wantBody: "{\"status\":\"OK\"}\n",
		},
		{
			name:     "empty body",
			mock:     nil,
			body:     ``,
			uid:      1,
			wantCode: http.StatusBadRequest,
			wantBody: "{\"status\":\"Error\",\"error\":\"empty request\"}\n",
		},
		{
			name:     "invalid email",
			mock:     nil,
			body:     `{"login":"newlogin", "email":"invalidemail", "password":"newpassword"}`,
			wantCode: http.StatusBadRequest,
			wantBody: "{\"status\":\"Error\",\"error\":\"field Email is not a valid email\"}\n",
		},
		{
			name:     "short login",
			mock:     nil,
			body:     `{"login":"abc", "email":"test@example.com", "password":"newpassword"}`,
			wantCode: http.StatusBadRequest,
			wantBody: "{\"status\":\"Error\",\"error\":\"field Login is out of its range\"}\n",
		},
		{
			name:     "long login",
			mock:     nil,
			body:     `{"login":"abcabcabcabcabcabc", "email":"test@example.com", "password":"newpassword"}`,
			wantCode: http.StatusBadRequest,
			wantBody: "{\"status\":\"Error\",\"error\":\"field Login is out of its range\"}\n",
		},
		{
			name:     "short password",
			mock:     nil,
			body:     `{"login":"login", "email":"test@example.com", "password":"pass"}`,
			wantCode: http.StatusBadRequest,
			wantBody: "{\"status\":\"Error\",\"error\":\"field Password is out of its range\"}\n",
		},
		{
			name:     "long password",
			mock:     nil,
			body:     `{"login":"login", "email":"test@example.com", "password":"newpasswordnewpasswordnewpassword"}`,
			wantCode: http.StatusBadRequest,
			wantBody: "{\"status\":\"Error\",\"error\":\"field Password is out of its range\"}\n",
		},
		{
			name: "invalid credentials",
			mock: mockUserEditor.EXPECT().
				EditUser(dto.EditUserDTO{
					Id:       1,
					Login:    "newlogin",
					Email:    "test@example.com",
					Password: "newpassword",
				}).
				Return(service.ErrInvalidCredentials).
				Times(1),
			body:     `{"login":"newlogin", "email":"test@example.com", "password":"newpassword"}`,
			uid:      1,
			wantCode: http.StatusUnauthorized,
			wantBody: "{\"status\":\"Error\",\"error\":\"invalid credentials\"}\n",
		},
		{
			name: "failed to update refresh token",
			mock: mockUserEditor.EXPECT().
				EditUser(dto.EditUserDTO{
					Id:       1,
					Login:    "newlogin",
					Email:    "test@example.com",
					Password: "newpassword",
				}).
				Return(errors.New("any")).
				Times(1),
			body:     `{"login":"newlogin", "email":"test@example.com", "password":"newpassword"}`,
			uid:      1,
			wantCode: http.StatusInternalServerError,
			wantBody: "{\"status\":\"Error\",\"error\":\"failed to update refresh token\"}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(method, URL, bytes.NewBufferString(tt.body))

			ctx := context.WithValue(r.Context(), "uid", int64(1))
			r = r.WithContext(ctx)

			handler := auth.Edit(log, mockUserEditor)
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
