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
	"regexp"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type EditRequest struct {
	Login    string `json:"login" validate:"required,min=4,max=15"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=30"`
}

type EditResponse struct {
	resp.Response
}

type UserEditor interface {
	EditUser(user dto.EditUserDTO) error
}

func Edit(log *slog.Logger, editor UserEditor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = packagePath + "edit.Edit"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req EditRequest

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

		loginMatch := regexp.MustCompile(`^[A-Za-z0-9]*$`)
		if !loginMatch.MatchString(req.Login) {
			log.Error("invalid login", slog.String("login", req.Login))
			responseError(w, r, resp.Error("invalid login"), http.StatusBadRequest)
			return
		}

		passwordMatch := regexp.MustCompile(`^[A-Za-z0-9\d@$!%*#?&]*$`)
		if !passwordMatch.MatchString(req.Password) {
			log.Error("invalid password", slog.String("password", req.Password))
			responseError(w, r, resp.Error("invalid password"), http.StatusBadRequest)
			return
		}

		if err = validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request", sl.Err(err))
			responseError(w, r, resp.ValidationError(validateErr), http.StatusBadRequest)
			return
		}

		id := r.Context().Value("uid").(int64)

		err = editor.EditUser(dto.EditUserDTO{
			Id:       id,
			Login:    req.Login,
			Email:    req.Email,
			Password: req.Password,
		})
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

		responseOK(w, r, resp.OK())
	}
}
