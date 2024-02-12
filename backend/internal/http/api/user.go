package api

import (
	"cloud-render/internal/dto"
	resp "cloud-render/internal/lib/response"
	"cloud-render/internal/lib/sl"
	"cloud-render/internal/service"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

type UserResposne struct {
	resp.Response
	Login      string     `json:"login"`
	Email      string     `json:"email"`
	ExpireDate *time.Time `json:"expirationDate"`
}

type UserInfoProvider interface {
	GetExpireDateWithUserInfo(uid int64) (*dto.UserInfoDTO, error)
}

func User(log *slog.Logger, userInfoProvider UserInfoProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = packagePath + "user.User"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		id := r.Context().Value("uid").(int64)

		userDTO, err := userInfoProvider.GetExpireDateWithUserInfo(id)
		if err != nil {
			if errors.Is(err, service.ErrExternalError) {
				log.Error("failed to get user info", sl.Err(err))
				responseError(w, r, resp.Error("failed to fetch user info"), http.StatusBadRequest)
				return
			}
			if errors.Is(err, service.ErrFailedToGetSubscription) {
				log.Error("failed to get subcription info", sl.Err(err))
				responseError(w, r, resp.Error("failed to get subcription info"), http.StatusInternalServerError)
				return
			}
			log.Error("failed to get user data", sl.Err(err))
			responseError(w, r, resp.Error("failed to get user data"), http.StatusInternalServerError)
			return
		}

		responseOK(w, r, UserResposne{
			Response:   resp.OK(),
			Login:      userDTO.Login,
			Email:      userDTO.Email,
			ExpireDate: userDTO.ExpirationDate,
		})
	}
}
