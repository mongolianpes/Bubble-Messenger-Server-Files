package uploader

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"time"

	"files/internal/rdb"
	"files/internal/storage"
)

type FilesRequest struct {
	File []byte
}

type UploaderService struct {
	temp    rdb.TempGetter
	storage storage.Storage
}

type Uploader interface {
	StartFilesHTTPServer(port, storagePath string)
}

func NewUploaderService() *UploaderService {
	return &UploaderService{}
}

func (s *UploaderService) StartFilesHTTPServer(port, storagePath string) {
	http.HandleFunc("/", s.FilesHandler)

	err := http.ListenAndServe(port, nil)
	if err != nil {
		slog.Error("Error start HTTP uploader", "error", err)
	}
}

func (s *UploaderService) FilesHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	storagePath, err := s.temp.GetStoragePathByTempPath(ctx, r.URL.Path)
	if err != nil {
		slog.Error("Error get storage path", "error", err)
		w.Write([]byte("Error get storage path"))
		return
	}

	switch storagePath {
	case "save_avatar":
		body, err := io.ReadAll(r.Body)
		if err != nil {
			slog.Error("Error read request body", "error", err)
			w.Write([]byte("Error read request body"))
			return
		}

		req := &FilesRequest{}
		if err = json.Unmarshal(body, &req); err != nil {
			w.Write([]byte("Error unmarshal request body"))
			return
		}

		storagePath, err = s.storage.SaveFile(req.File, true)
		if err != nil {
			slog.Error("Error save file", "error", err)
			w.Write([]byte("Error save file"))
			return
		}
	case "save_file":
		body, err := io.ReadAll(r.Body)
		if err != nil {
			slog.Error("Error read request body", "error", err)
			w.Write([]byte("Error read request body"))
			return
		}

		req := &FilesRequest{}
		if err = json.Unmarshal(body, &req); err != nil {
			w.Write([]byte("Error unmarshal request body"))
			return
		}

		storagePath, err = s.storage.SaveFile(req.File, false)
		if err != nil {
			slog.Error("Error save file", "error", err)
			w.Write([]byte("Error save file"))
			return
		}
	default:
		file, err := s.storage.GetFile(storagePath)
		if err != nil {
			slog.Error("Cant read file", "error", err)
			w.Write([]byte("Error read file"))
			return
		}

		w.Write(file)
		return
	}
}
