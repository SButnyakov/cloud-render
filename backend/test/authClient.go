package test

import (
	"bytes"
	authHandler "cloud-render/internal/http/auth"
	"cloud-render/internal/lib/config"
	"encoding/json"
	"fmt"
)

type AuthClient struct {
	AccessToken  string
	RefreshToken string
	Config       *config.Config
}

func NewAuthClient(config *config.Config) *AuthClient {
	return &AuthClient{Config: config}
}

func (c *AuthClient) SignUp(login, email, password string) (*ResponseParams, error) {
	req := authHandler.SignUpRequest{
		Login:    login,
		Email:    email,
		Password: password,
	}

	jsonBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	return MakeRequest(RequestParams{
		Method: "POST",
		URL: fmt.Sprintf("http://%s:%d%s",
			c.Config.HTTPServer.Host,
			c.Config.HTTPServer.Port,
			c.Config.Paths.SignUp),
		Body: bytes.NewReader(jsonBody),
	})
}

func (c *AuthClient) SignIn(loginOrEmail, password string) (*ResponseParams, error) {
	req := authHandler.SignInRequest{
		LoginOrEmail: loginOrEmail,
		Password:     password,
	}

	jsonBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	params, err := MakeRequest(RequestParams{
		Method: "POST",
		URL: fmt.Sprintf("http://%s:%d%s",
			c.Config.HTTPServer.Host,
			c.Config.HTTPServer.Port,
			c.Config.Paths.SignIn),
		Body: bytes.NewReader(jsonBody),
	})
	if err != nil {
		return nil, err
	}

	res := authHandler.SignInResponse{}

	err = json.Unmarshal(params.Body, &res)
	if err != nil {
		return nil, err
	}

	c.AccessToken = res.AccessToken
	c.RefreshToken = res.RefreshToken

	return params, nil
}

func (c *AuthClient) Refresh() (*ResponseParams, error) {
	req := authHandler.RefreshRequest{
		RefreshToken: c.RefreshToken,
	}

	jsonBody, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	params, err := MakeRequest(RequestParams{
		Method: "POST",
		URL: fmt.Sprintf("http://%s:%d%s",
			c.Config.HTTPServer.Host,
			c.Config.HTTPServer.Port,
			c.Config.Paths.Refresh),
		Body: bytes.NewReader(jsonBody),
	})
	if err != nil {
		return nil, err
	}

	res := authHandler.RefreshResponse{}

	err = json.Unmarshal(params.Body, &res)
	if err != nil {
		return nil, err
	}

	c.AccessToken = res.AccessToken
	c.RefreshToken = res.RefreshToken

	return params, nil
}
