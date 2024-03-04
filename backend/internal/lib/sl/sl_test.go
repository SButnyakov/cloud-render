package sl_test

import (
	"cloud-render/internal/lib/sl"
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestSl_SetupLogger(t *testing.T) {
	localLogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	devLogger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	prodLogger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	assert.Equal(t, localLogger, sl.SetupLogger("local"))
	assert.Equal(t, devLogger, sl.SetupLogger("dev"))
	assert.Equal(t, prodLogger, sl.SetupLogger("prod"))
}

func TestSl_Err(t *testing.T) {
	errMsg := "error message"

	assert.Equal(t, sl.Err(errors.New(errMsg)), slog.Attr{
		Key:   "error",
		Value: slog.StringValue(errMsg),
	})
}
