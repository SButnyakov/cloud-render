package auth

import (
	"cloud-render/internal/dto"
	resp "cloud-render/internal/lib/response"
	"cloud-render/internal/lib/sl"
	"cloud-render/internal/service"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type InfoResponse struct {
	resp.Response
	Login string `json:"login"`
	Email string `json:"email"`
}

type UserProvider interface {
	GetUser(id int64) (*dto.GetUserDTO, error)
}

func Info(log *slog.Logger, userProvider UserProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = packagePath + "info.Info"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		id := chi.URLParam(r, "id")
		if id == "" {
			log.Error("id not provided")
			responseError(w, r, resp.Error("invalid id"), http.StatusBadRequest)
			return
		}

		idInt64, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			log.Error("invalid id", sl.Err(err))
			responseError(w, r, resp.Error("invalid id"), http.StatusBadRequest)
			return
		}

		userDTO, err := userProvider.GetUser(idInt64)
		if err != nil {
			if errors.Is(err, service.ErrUserNotFound) {
				log.Error("user not found")
				responseError(w, r, resp.Error("user not found"), http.StatusNotFound)
				return
			}
		}

		responseOK(w, r, InfoResponse{
			Response: resp.OK(),
			Login:    userDTO.Login,
			Email:    userDTO.Email,
		})
	}
}
