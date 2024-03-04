package auth_test

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"cloud-render/internal/http/middleware/auth"
	mocks "cloud-render/internal/mocks/lib"
)

func TestAuthMiddleware(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockTokenManager := mocks.NewMockTokenManager(mockCtrl)

	middleware := auth.New(log, mockTokenManager)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	})

	tests := []struct {
		name         string
		authHeader   string
		token        string
		expiresAt    int64
		expectedCode int
		expectedBody string
		mocks        *gomock.Call
	}{
		{
			name:         "Valid token",
			authHeader:   "Bearer valid_token",
			token:        "valid_token",
			expiresAt:    time.Now().Add(time.Hour).Unix(),
			expectedCode: http.StatusOK,
			mocks: mockTokenManager.EXPECT().
				Parse("valid_token").
				Return(&jwt.StandardClaims{Subject: "1", ExpiresAt: time.Now().Add(time.Hour).Unix()}, nil),
		},

		{
			name:         "Empty authorization header",
			authHeader:   "",
			expectedCode: http.StatusUnauthorized,
			expectedBody: "{\"status\":\"Error\",\"error\":\"empty authorization header\"}\n",
		},
		{
			name:         "invalid authorization header",
			authHeader:   "invalid",
			expectedCode: http.StatusUnauthorized,
			expectedBody: "{\"status\":\"Error\",\"error\":\"invalid authorization header\"}\n",
		},
		{
			name:         "empty authorization token",
			authHeader:   "Bearer ",
			expectedCode: http.StatusUnauthorized,
			expectedBody: "{\"status\":\"Error\",\"error\":\"empty authorization token\"}\n",
		},
		{
			name:         "expired token",
			authHeader:   "Bearer expired_token",
			expectedCode: http.StatusUnauthorized,
			mocks: mockTokenManager.EXPECT().
				Parse("expired_token").
				Return(&jwt.StandardClaims{Subject: "1", ExpiresAt: time.Now().Add(-time.Hour).Unix()}, nil),
			expectedBody: "{\"status\":\"Error\",\"error\":\"token has expired\"}\n",
		},
		{
			name:         "invalid payload",
			authHeader:   "Bearer invalid_token",
			expectedCode: http.StatusUnauthorized,
			mocks: mockTokenManager.EXPECT().
				Parse("invalid_token").
				Return(&jwt.StandardClaims{Subject: "not int", ExpiresAt: time.Now().Add(time.Hour).Unix()}, nil),
			expectedBody: "{\"status\":\"Error\",\"error\":\"invalid payload\"}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			ctx := context.WithValue(req.Context(), "uid", int64(1))
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()

			middleware(handler).ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			assert.Equal(t, tt.expectedBody, w.Body.String())
		})
	}
}
