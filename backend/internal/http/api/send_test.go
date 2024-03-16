package api_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"cloud-render/internal/http/api"
	mocks "cloud-render/internal/mocks/service/api"
)

func TestSendHandler_Send(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockOrderCreator := mocks.NewMockOrderCreator(mockCtrl)
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	sendHandler := api.Send(log, mockOrderCreator)

	t.Run("Success", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		formFile, _ := writer.CreateFormFile("uploadfile", "testfile.txt")
		io.WriteString(formFile, "Test file content")
		writer.WriteField("format", "jpeg")
		writer.WriteField("resolution", "1024x768")
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/send", body)
		ctx := context.WithValue(req.Context(), "uid", int64(1))
		req = req.WithContext(ctx)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		mockOrderCreator.EXPECT().CreateOrder(gomock.Any()).Return(nil)

		recorder := httptest.NewRecorder()

		sendHandler.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusCreated, recorder.Code)
		assert.Equal(t, "{\"status\":\"OK\"}\n", recorder.Body.String())
	})

	t.Run("Failed to parse form", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/send", nil)
		ctx := context.WithValue(req.Context(), "uid", int64(1))
		req = req.WithContext(ctx)

		recorder := httptest.NewRecorder()
		sendHandler.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Equal(t, "{\"status\":\"Error\",\"error\":\"failed to parse form\"}\n", recorder.Body.String())
	})

	t.Run("Failed to get file from form", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		writer.WriteField("format", "jpeg")
		writer.WriteField("resolution", "1024x768")
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/send", body)
		ctx := context.WithValue(req.Context(), "uid", int64(1))
		req = req.WithContext(ctx)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		recorder := httptest.NewRecorder()
		sendHandler.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Equal(t, "{\"status\":\"Error\",\"error\":\"failed to get file from form\"}\n", recorder.Body.String())
	})

	t.Run("Failed to create new order", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		formFile, _ := writer.CreateFormFile("uploadfile", "testfile.txt")
		io.WriteString(formFile, "Test file content")
		writer.WriteField("format", "jpeg")
		writer.WriteField("resolution", "1024x768")
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/send", body)
		ctx := context.WithValue(req.Context(), "uid", int64(1))
		req = req.WithContext(ctx)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		mockOrderCreator.EXPECT().CreateOrder(gomock.Any()).Return(errors.New("failed to create order"))

		recorder := httptest.NewRecorder()
		sendHandler.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		assert.Equal(t, "{\"status\":\"Error\",\"error\":\"failed to create new order\"}\n", recorder.Body.String())
	})
}
