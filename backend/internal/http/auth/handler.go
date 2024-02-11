package auth

import (
	resp "cloud-render/internal/lib/response"
	"net/http"

	"github.com/go-chi/render"
)

const (
	packagePath = "http.auth.signin."
)

func responseError(w http.ResponseWriter, r *http.Request, response resp.Response, status int) {
	w.WriteHeader(status)
	render.JSON(w, r, response)
}

func responseOK(w http.ResponseWriter, r *http.Request, v interface{}) {
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, Response{
		Response: resp.OK(),
	})
}
