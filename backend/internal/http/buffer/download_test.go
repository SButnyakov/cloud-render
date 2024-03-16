package buffer_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"cloud-render/internal/http/buffer"
	"log/slog"

	"github.com/stretchr/testify/assert"
)

func TestDownloadHandler(t *testing.T) {
	tempDir := t.TempDir()

	fileContent := []byte("This is a test file.")
	filePath := filepath.Join(tempDir, "test_file.txt")
	err := os.WriteFile(filePath, fileContent, 0666)
	if err != nil {
		t.Fatal("failed to create test file:", err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	r := httptest.NewRequest(http.MethodGet, "/download/123/test_file.txt", nil)
	w := httptest.NewRecorder()

	handler := buffer.Download(logger, tempDir)
	handler(w, r)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode, "status code should be OK")
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		t.Fatal("failed to read response body:", err)
	}
	assert.Equal(t, fileContent, buf.Bytes(), "file content should match")

	r = httptest.NewRequest(http.MethodGet, "/download/123/non_existing_file.txt", nil)
	w = httptest.NewRecorder()

	handler(w, r)

	resp = w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "status code should be BadRequest")
}
