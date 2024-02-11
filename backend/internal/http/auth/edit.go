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
	"strconv"

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
	EditUer(user dto.EditUserDTO) error
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

		if err = validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request", sl.Err(err))
			responseError(w, r, resp.ValidationError(validateErr), http.StatusBadRequest)
			return
		}

		id := r.Context().Value("uid").(string)
		idInt64, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			log.Error("invalid id", sl.Err(err))
			responseError(w, r, resp.Error("invalid id"), http.StatusBadRequest)
			return
		}

		err = editor.EditUer(dto.EditUserDTO{
			Id:       idInt64,
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
