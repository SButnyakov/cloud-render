package buffer

import (
	"cloud-render/internal/lib/config"
	resp "cloud-render/internal/lib/response"
	"cloud-render/internal/lib/sl"
	"cloud-render/internal/models"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type RequestResponse struct {
	resp.Response
	Format       string `json:"format"`
	Resolution   string `json:"resolution"`
	DownloadLink string `json:"download_link,omitempty"`
}

func Request(log *slog.Logger, client *redis.Client, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = packagePath + "request.Request"

		data, err := client.BLPop(context.Background(), time.Second, cfg.Redis.PriorityQueueName).Result()
		if err != nil && !errors.Is(err, redis.Nil) {
			log.Error("reading redis priority queue fail", sl.Err(err))
		}
		if err != nil {
			data, err = client.BLPop(context.Background(), time.Second, cfg.Redis.QueueName).Result()
		}
		if err != nil {
			if errors.Is(err, redis.Nil) {
				log.Info("empty queue's")
				responseEmpty(w, r)
				return
			}
			log.Error("reading redis queue fail")
			responseError(w, r, resp.Error("reading orders list failed"), http.StatusInternalServerError)
			return
		}

		var newOrder models.RedisData

		b := []byte(data[1])
		err = json.Unmarshal(b, &newOrder)
		if err != nil {
			log.Error("failed to unmarshal new order", sl.Err(err))
		}

		fmt.Println(newOrder)

		pathList := strings.Split(newOrder.SavePath, "/")
		listLength := len(pathList)

		downloadLink := fmt.Sprintf("http://%s:%d/%s/blend/download/%s", cfg.HTTPServer.Host, cfg.HTTPServer.Port, pathList[listLength-2], pathList[listLength-1])

		responseOK(w, r, RequestResponse{
			Response:     resp.OK(),
			Format:       newOrder.Format,
			Resolution:   newOrder.Resolution,
			DownloadLink: downloadLink,
		})
	}
}
