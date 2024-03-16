package response_test

import (
	"cloud-render/internal/lib/response"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResponse_OK(t *testing.T) {
	resp := response.OK()

	assert.Equal(t, resp.Status, response.StatusOK)
	assert.Equal(t, resp.Error, "")
}

func TestResponse_Empty(t *testing.T) {
	resp := response.Empty()

	assert.Equal(t, resp.Status, response.StatusEmpty)
	assert.Equal(t, resp.Error, "")
}

func TestResponse_Error(t *testing.T) {
	resp := response.Error("error msg")

	assert.Equal(t, resp.Status, response.StatusError)
	assert.Equal(t, resp.Error, "error msg")
}
