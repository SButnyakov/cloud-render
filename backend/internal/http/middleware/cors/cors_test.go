package cors_test

import (
	"cloud-render/internal/http/middleware/cors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCors(t *testing.T) {
	corsHandler := cors.New()
	require.NotNil(t, corsHandler)
}
