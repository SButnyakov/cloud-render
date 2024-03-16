package buffer

import (
	resp "cloud-render/internal/lib/response"
	"net/http"

	"github.com/go-chi/render"
)

const (
	packagePath = "http.buffer."
)

func responseError(w http.ResponseWriter, r *http.Request, response resp.Response, status int) {
	w.WriteHeader(status)
	render.JSON(w, r, response)
}

func responseOK(w http.ResponseWriter, r *http.Request, v interface{}, status ...int) {
	responseStatus := http.StatusOK
	if len(status) > 0 {
		responseStatus = status[0]
	}
	w.WriteHeader(responseStatus)
	render.JSON(w, r, v)
}

func responseEmpty(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, resp.Empty())
}
