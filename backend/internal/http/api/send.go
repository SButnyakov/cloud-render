package api

import (
	"cloud-render/internal/dto"
	resp "cloud-render/internal/lib/response"
	"cloud-render/internal/lib/sl"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

type OrderCreator interface {
	CreateOrder(dto dto.CreateOrderDTO) error
}

func Send(log *slog.Logger, orderCreator OrderCreator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = packagePath + "send.Send"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			log.Error("failed to parse form", sl.Err(err))
			responseError(w, r, resp.Error("failed to parse form"), http.StatusBadRequest)
			return
		}

		format := r.FormValue("format")
		resolution := r.FormValue("resolution")

		file, header, err := r.FormFile("uploadfile")
		if err != nil {
			log.Error("failed to get file from form", sl.Err(err))
			responseError(w, r, resp.Error("failed to get file from form"), http.StatusBadRequest)
			return
		}
		defer file.Close()

		id := r.Context().Value("uid").(int64)

		log.Info("orderDTO", slog.Any("order", dto.CreateOrderDTO{
			UserId:     id,
			Format:     format,
			Resolution: resolution,
			File:       file,
			Header:     header,
		}))

		err = orderCreator.CreateOrder(dto.CreateOrderDTO{
			UserId:     id,
			Format:     format,
			Resolution: resolution,
			File:       file,
			Header:     header,
		})
		if err != nil {
			log.Error("failed to create new order", sl.Err(err))
			responseError(w, r, resp.Error("failed to create new order"), http.StatusInternalServerError)
			return
		}

		responseOK(w, r, resp.OK(), http.StatusCreated)
	}
}
