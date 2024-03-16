package auth

import (
	"cloud-render/internal/dto"
	resp "cloud-render/internal/lib/response"
	"cloud-render/internal/lib/sl"
	"cloud-render/internal/service"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type RefreshResponse struct {
	resp.Response
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserReauthorizer interface {
	ReauthUser(userDTO dto.ReAuthUserDTO) (*dto.ReAuthUserDTO, error)
}

func Refresh(log *slog.Logger, reauthorizer UserReauthorizer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = packagePath + "refresh.Refresh"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req RefreshRequest

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty", sl.Err(err))
			responseError(w, r, resp.Error("empty request"), http.StatusBadRequest)
			return
		}
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			responseError(w, r, resp.Error("failed to decode request"), http.StatusBadRequest)
			return
		}

		if err = validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request", sl.Err(err))
			responseError(w, r, resp.ValidationError(validateErr), http.StatusBadRequest)
			return
		}

		responseDTO, err := reauthorizer.ReauthUser(dto.ReAuthUserDTO{RefreshToken: req.RefreshToken})
		if err != nil {
			if errors.Is(err, service.ErrInvalidCredentials) {
				log.Debug("invalid credentials")
				responseError(w, r, resp.Error("invalid credentials"), http.StatusUnauthorized)
				return
			}
			log.Error("failed to authorize user", sl.Err(err))
			responseError(w, r, resp.Error("failed to update refresh token"), http.StatusInternalServerError)
			return
		}

		responseOK(w, r, RefreshResponse{
			Response:     resp.OK(),
			AccessToken:  responseDTO.AccessToken,
			RefreshToken: responseDTO.RefreshToken,
		})
	}
}
