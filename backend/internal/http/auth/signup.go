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

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type SignUpRequest struct {
	Login    string `json:"login" validate:"required,min=4,max=15"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=30"`
}

type SignUpResponse struct {
	resp.Response
}

type UserCreator interface {
	CreateUser(dto.CreateUserDTO) error
}

func SignUp(log *slog.Logger, userCreator UserCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = packagePath + "signup.SignUp"

		var req SignUpRequest

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

		err = userCreator.CreateUser(dto.CreateUserDTO{
			Login:    req.Login,
			Email:    req.Email,
			Password: req.Password,
		})
		if errors.Is(err, service.ErrUserAlreadyExists) {
			responseError(w, r, resp.Error("user with the same login or email already exists"), http.StatusBadRequest)
			return
		}
		if err != nil {
			log.Error("failed to create new user", sl.Err(err))
			responseError(w, r, resp.Error("server-side registration fail"), http.StatusBadRequest)
			return
		}

		responseOK(w, r, resp.OK(), http.StatusCreated)
	}
}
