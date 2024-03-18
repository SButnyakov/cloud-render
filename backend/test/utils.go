package test

import (
	"io"
	"net/http"
)

type HeadersMap map[string]string

type RequestParams struct {
	Method  string
	URL     string
	Headers HeadersMap
	Body    io.Reader
}

type ResponseParams struct {
	Code int
	Body []byte
}

func MakeRequest(params RequestParams) (*ResponseParams, error) {
	req, err := http.NewRequest(params.Method, params.URL, params.Body)
	if err != nil {
		return nil, err
	}

	for k, v := range params.Headers {
		req.Header.Set(k, v)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return &ResponseParams{
		Code: res.StatusCode,
		Body: resBody,
	}, nil
}
