package api

import (
	"bytes"
	"cloud-render/internal/lib/config"
	"cloud-render/test"
	"fmt"
	"io"
	"mime/multipart"
)

type APITestClient struct {
	*test.AuthClient
	Config *config.Config
}

func NewAPITestClient(config *config.Config, authClient *test.AuthClient) *APITestClient {
	return &APITestClient{
		AuthClient: authClient,
		Config:     config,
	}
}

func (c *APITestClient) User() (*test.ResponseParams, error) {
	return test.MakeRequest(test.RequestParams{
		Method: "GET",
		URL: fmt.Sprintf("http://%s:%d%s",
			c.Config.HTTPServer.Host,
			c.Config.HTTPServer.Port,
			c.Config.Paths.User),
		Headers: test.HeadersMap{
			"Authorization": fmt.Sprintf("Bearer %s", c.AccessToken),
		},
	})
}

func (c *APITestClient) Orders() (*test.ResponseParams, error) {
	return test.MakeRequest(test.RequestParams{
		Method: "GET",
		URL: fmt.Sprintf("http://%s:%d%s",
			c.Config.HTTPServer.Host,
			c.Config.HTTPServer.Port,
			c.Config.Paths.Orders.Root),
		Headers: test.HeadersMap{
			"Authorization": fmt.Sprintf("Bearer %s", c.AccessToken),
		},
	})
}

func (c *APITestClient) Send(filename, format, resolution string) (*test.ResponseParams, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	formFile, _ := writer.CreateFormFile("uploadfile", filename)
	io.WriteString(formFile, "Test file content")
	writer.WriteField("format", format)
	writer.WriteField("resolution", resolution)
	writer.Close()

	return test.MakeRequest(test.RequestParams{
		Method: "POST",
		URL: fmt.Sprintf("http://%s:%d%s",
			c.Config.HTTPServer.Host,
			c.Config.HTTPServer.Port,
			c.Config.Paths.Send),
		Headers: test.HeadersMap{
			"Authorization": fmt.Sprintf("Bearer %s", c.AccessToken),
			"Content-Type":  writer.FormDataContentType(),
		},
		Body: body,
	})
}

func (c *APITestClient) SoftDelete(orderId string) (*test.ResponseParams, error) {
	return test.MakeRequest(test.RequestParams{
		Method: "POST",
		URL: fmt.Sprintf("http://%s:%d%s/%s/delete",
			c.Config.HTTPServer.Host,
			c.Config.HTTPServer.Port,
			c.Config.Paths.Orders.Root,
			orderId),
		Headers: test.HeadersMap{
			"Authorization": fmt.Sprintf("Bearer %s", c.AccessToken),
		},
	})
}

func (c *APITestClient) Subscribe() (*test.ResponseParams, error) {
	return test.MakeRequest(test.RequestParams{
		Method: "POST",
		URL: fmt.Sprintf("http://%s:%d%s",
			c.Config.HTTPServer.Host,
			c.Config.HTTPServer.Port,
			c.Config.Paths.Subscribe),
		Headers: test.HeadersMap{
			"Authorization": fmt.Sprintf("Bearer %s", c.AccessToken),
		},
	})
}
