package buffer

import (
	"cloud-render/internal/lib/sl"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Download(log *slog.Logger, downloadPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = packagePath + "download.Download"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		uid := chi.URLParam(r, "uid")
		fileName := path.Base(r.URL.Path)
		filePath := filepath.Join(downloadPath, uid, fileName)

		file, err := os.Open(filePath)
		if err != nil {
			log.Error("file not found", sl.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer file.Close()

		tempBuffer := make([]byte, 512)
		_, err = file.Read(tempBuffer)
		if err != nil {
			log.Error("failed to read file to buffer", sl.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		fileContentType := http.DetectContentType(tempBuffer)

		fileStat, err := file.Stat()
		if err != nil {
			log.Error("failed to get file stat", sl.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		fileSize := strconv.FormatInt(fileStat.Size(), 10)

		w.Header().Set("Content-Type", fileContentType+";"+fileName)
		w.Header().Set("Content-Length", fileSize)

		_, err = file.Seek(0, 0)
		if err != nil {
			log.Error("failed to reset offset of file")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, err = io.Copy(w, file)
		if err != nil {
			log.Error("failed to send the file to client")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
