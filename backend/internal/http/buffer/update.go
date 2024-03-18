package buffer

import (
	"cloud-render/internal/dto"
	resp "cloud-render/internal/lib/response"
	"cloud-render/internal/lib/sl"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type OrderStatusUpdater interface {
	UpdateOrderStatus(dto dto.UpdateOrderStatusDTO) error
}

func Update(log *slog.Logger, updater OrderStatusUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = packagePath + "update.Update"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		id := chi.URLParam(r, "id")
		idInt64, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			log.Error("invalid order id", sl.Err(err))
			responseError(w, r, resp.Error("invalid order id"), http.StatusBadRequest)
			return
		}

		status := strings.ReplaceAll(chi.URLParam(r, "status"), "-", " ")

		err = updater.UpdateOrderStatus(dto.UpdateOrderStatusDTO{
			OrderId: idInt64,
			Status:  status,
		})
		if err != nil {
			log.Error("failed to update status", sl.Err(err))
			responseError(w, r, resp.Error("failed to update status"), http.StatusInternalServerError)
			return
		}

		responseOK(w, r, resp.OK())
	}
}
