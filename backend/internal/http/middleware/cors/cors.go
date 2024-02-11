package cors

import (
	"github.com/go-chi/cors"
	"net/http"
)

func New() func(http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowOriginFunc:  AllowOriginFunc,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
}

func AllowOriginFunc(r *http.Request, origin string) bool {
	return origin == "http://localhost:3000" || origin == "http://localhost:8080"
}
