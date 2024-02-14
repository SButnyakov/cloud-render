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

		id := chi.URLParam(r, "uid")
		idInt64, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			log.Error("invalid user id", sl.Err(err))
			responseError(w, r, resp.Error("invalid user id"), http.StatusBadRequest)
			return
		}

		fileName := chi.URLParam(r, "filename")
		status := strings.ReplaceAll(chi.URLParam(r, "status"), "-", " ")

		err = updater.UpdateOrderStatus(dto.UpdateOrderStatusDTO{
			UserId:      idInt64,
			StoringName: fileName,
			Status:      status,
		})

		responseOK(w, r, resp.OK())
	}
}
