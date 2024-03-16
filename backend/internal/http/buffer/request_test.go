package buffer_test

import (
	"cloud-render/internal/http/buffer"
	"cloud-render/internal/lib/config"
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
)

func TestRequestHandler_Request(t *testing.T) {
	s, err := miniredis.Run()
	if err != nil {
		t.Fatalf("failed to start miniredis: %v", err)
	}
	defer s.Close()

	client := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	mockCfg := &config.Config{
		Redis: config.Redis{
			PriorityQueueName: "priority_queue",
			QueueName:         "queue",
		},
		HTTPServer: config.HTTPServer{
			Host: "localhost",
			Port: 8080,
		},
	}

	tests := []struct {
		name          string
		redisPrep     *redis.IntCmd
		mockRedis     *redis.Client
		expectedCode  int
		expectedBody  string
		expectedError error
	}{
		{
			name:          "successful priority request",
			redisPrep:     client.RPush(context.Background(), mockCfg.Redis.PriorityQueueName, string(`{"format":"jpeg","resolution":"1080p","save_path":"input/1/123123"}`)),
			mockRedis:     client,
			expectedCode:  http.StatusOK,
			expectedBody:  "{\"status\":\"OK\",\"format\":\"jpeg\",\"resolution\":\"1080p\",\"download_link\":\"http://localhost:8080/1/blend/download/123123\"}\n",
			expectedError: nil,
		},
		{
			name:          "successful request",
			redisPrep:     client.RPush(context.Background(), mockCfg.Redis.QueueName, string(`{"format":"jpeg","resolution":"1080p","save_path":"input/1/123123"}`)),
			mockRedis:     client,
			expectedCode:  http.StatusOK,
			expectedBody:  "{\"status\":\"OK\",\"format\":\"jpeg\",\"resolution\":\"1080p\",\"download_link\":\"http://localhost:8080/1/blend/download/123123\"}\n",
			expectedError: nil,
		},
		{
			name:          "empty queue",
			mockRedis:     client,
			expectedCode:  http.StatusOK,
			expectedBody:  "{\"status\":\"Empty\"}\n",
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
			handler := buffer.Request(log, tt.mockRedis, mockCfg)

			req := httptest.NewRequest("GET", "/request", nil)
			w := httptest.NewRecorder()

			handler(w, req)

			resp := w.Result()
			if resp.StatusCode != tt.expectedCode {
				t.Errorf("Expected status code %d, got %d", tt.expectedCode, resp.StatusCode)
			}

			body := w.Body.String()
			if body != tt.expectedBody {
				t.Errorf("Expected body %q, got %q", tt.expectedBody, body)
			}
		})
	}
}
