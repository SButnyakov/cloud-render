package buffer

import (
	"cloud-render/internal/lib/config"
	"cloud-render/test"
	"fmt"
)

type BufferTestClient struct {
	Config *config.Config
}

func NewBufferTestClient(config *config.Config) *BufferTestClient {
	return &BufferTestClient{
		Config: config,
	}
}

func (c *BufferTestClient) Update(uid, orderId, status string) (*test.ResponseParams, error) {
	return test.MakeRequest(test.RequestParams{
		Method: "PUT",
		URL: fmt.Sprintf("http://%s:%d/%s/blend/update/%s/%s",
			c.Config.HTTPServer.Host,
			c.Config.HTTPServer.Port,
			uid,
			orderId,
			status,
		),
	})
}

func (c *BufferTestClient) Request() (*test.ResponseParams, error) {
	return test.MakeRequest(test.RequestParams{
		Method: "GET",
		URL: fmt.Sprintf("http://%s:%d%s",
			c.Config.HTTPServer.Host,
			c.Config.HTTPServer.Port,
			c.Config.Paths.Request,
		),
	})
}
