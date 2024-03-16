package api

import (
	"cloud-render/internal/lib/config"
	resp "cloud-render/internal/lib/response"
	"cloud-render/internal/lib/sl"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

type UserSubscriber interface {
	SubscribeUser(int64) error
}

func Subscribe(log *slog.Logger, cfg *config.Config, subscriber UserSubscriber) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn := packagePath + "subscribe.Susbcribe"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		id := r.Context().Value("uid").(int64)

		err := subscriber.SubscribeUser(id)
		if err != nil {
			log.Error("failed to subscribe user", sl.Err(err))
			responseError(w, r, resp.Error("failed to subscribe"), http.StatusInternalServerError)
			return
		}

		responseOK(w, r, resp.OK())
	}
}
