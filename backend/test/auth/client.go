package auth

import (
	"bytes"
	authHandler "cloud-render/internal/http/auth"
	"cloud-render/internal/lib/config"
	"cloud-render/test"
	"encoding/json"
	"fmt"
)

type AuthTestClient struct {
	*test.AuthClient
	Config *config.Config
}

func NewAuthTestClient(config *config.Config, authClient *test.AuthClient) *AuthTestClient {
	return &AuthTestClient{
		AuthClient: authClient,
		Config:     config,
	}
}

func (c *AuthTestClient) Edit(login, email, password string) (*test.ResponseParams, error) {
	req := authHandler.EditRequest{
		Login:    login,
		Email:    email,
		Password: password,
	}

	jsonBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	return test.MakeRequest(test.RequestParams{
		Method: "PUT",
		URL: fmt.Sprintf("http://%s:%d%s",
			c.Config.HTTPServer.Host,
			c.Config.HTTPServer.Port,
			c.Config.Paths.Edit),
		Headers: test.HeadersMap{
			"Authorization": fmt.Sprintf("Bearer %s", c.AccessToken),
		},
		Body: bytes.NewReader(jsonBody),
	})
}
