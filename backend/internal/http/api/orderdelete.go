package api

import (
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

type OneOrderSoftDelter interface {
	SoftDeleteOneOrder(id int64) error
}

func DeleteOrder(log *slog.Logger, orderDeleter OneOrderSoftDelter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = packagePath + "orderdelete.DeleteOrder"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		id := chi.URLParam(r, "id")
		if id == "" {
			responseError(w, r, resp.Error("empty id param"), http.StatusBadRequest)
			return
		}

		int64Id, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			responseError(w, r, resp.Error("invalid id param"), http.StatusBadRequest)
			return
		}

		err = orderDeleter.SoftDeleteOneOrder(int64Id)
		if err != nil {
			if errors.Is(err, service.ErrOrderNotFound) {
				log.Error("order not found", slog.Int64(id, int64Id))
				responseError(w, r, resp.Error("order not found"), http.StatusNotFound)
				return
			}
			log.Error("failed to delete order", sl.Err(err))
			responseError(w, r, resp.Error("failed to delete order"), http.StatusInternalServerError)
			return
		}

		responseOK(w, r, resp.OK())
		return
	}
}
