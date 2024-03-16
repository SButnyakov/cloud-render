package buffer

import (
	"cloud-render/internal/dto"
	resp "cloud-render/internal/lib/response"
	"cloud-render/internal/lib/sl"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type OrdersImageUpdater interface {
	UpdateOrderImage(dto dto.UpdateOrderImageDTO) error
}

func Upload(log *slog.Logger, provider OrdersImageUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = packagePath + "upload.Upload"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			log.Error("failed to parse form", sl.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		file, header, err := r.FormFile("uploadfile")
		if err != nil {
			log.Error("failed to get file from form", sl.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer file.Close()

		uid := chi.URLParam(r, "uid")

		err = provider.UpdateOrderImage(dto.UpdateOrderImageDTO{
			UserId: uid,
			File:   file,
			Header: header,
		})
		if err != nil {
			log.Error("failed to upload image", sl.Err(err))
			responseError(w, r, resp.Error("failed to upload image"), http.StatusInternalServerError)
			return
		}

		responseOK(w, r, resp.OK())
	}
}
