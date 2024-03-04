package api

import (
	"cloud-render/internal/dto"
	resp "cloud-render/internal/lib/response"
	"cloud-render/internal/lib/sl"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

type GetManyOrdersResponse struct {
	resp.Response
	Orders []orderResponse `json:"orders"`
}

type orderResponse struct {
	Id           int64  `json:"id"`
	FileName     string `json:"filename"`
	CreationDate string `json:"date"`
	Status       string `json:"status"`
	DownloadLink string `json:"downloadLink,omitempty"`
}

type ManyOrdersProvider interface {
	GetManyOrders(id int64) ([]dto.GetOrderDTO, error)
}

func Orders(log *slog.Logger, orderProvider ManyOrdersProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = packagePath + "orders.Orders"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		id := r.Context().Value("uid").(int64)

		orders, err := orderProvider.GetManyOrders(id)
		if err != nil {
			log.Error("failed to get orders", sl.Err(err))
			responseError(w, r, resp.Error("failed to get orders"), http.StatusInternalServerError)
			return
		}

		responseOrders := make([]orderResponse, len(orders))
		for i, v := range orders {
			responseOrders[i] = orderResponse{
				Id:           v.Id,
				FileName:     v.Filename,
				CreationDate: v.Date.Format("02-01-2006"),
				Status:       v.OrderStatus,
				DownloadLink: v.DownloadLink,
			}
		}

		responseOK(w, r, GetManyOrdersResponse{
			Response: resp.OK(),
			Orders:   responseOrders,
		})
	}
}
