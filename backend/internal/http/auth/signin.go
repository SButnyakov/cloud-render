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

type SignInRequest struct {
	LoginOrEmail string `json:"login_or_email" validate:"required"`
	Password     string `json:"password" validate:"required"`
}

type SignInResponse struct {
	resp.Response
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserAuthorizer interface {
	AuthUser(userDTO dto.AuthUserDTO) (*dto.AuthUserDTO, error)
}

func SignIn(log *slog.Logger, userAuthorizer UserAuthorizer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = packagePath + "signin.SignIn"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req SignInRequest

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")
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

		respDTO, err := userAuthorizer.AuthUser(dto.AuthUserDTO{
			LoginOrEmail: req.LoginOrEmail,
			Password:     req.Password,
		})
		if err != nil {
			if errors.Is(err, service.ErrInvalidCredentials) {
				log.Debug("invalid credentials", sl.Err(err))
				responseError(w, r, resp.Error("invalid credentials"), http.StatusBadRequest)
				return
			}
			log.Error("server-side authorization fail", sl.Err(err))
			responseError(w, r, resp.Error("server-side authorization failed"), http.StatusInternalServerError)
			return
		}

		responseOK(w, r, SignInResponse{
			Response:     resp.OK(),
			AccessToken:  respDTO.AccessToken,
			RefreshToken: respDTO.RefreshToken,
		})
	}
}
