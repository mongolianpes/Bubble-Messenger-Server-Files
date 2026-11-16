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
	"files/internal/waiter"
)

type FilesRequest struct {
	File          []byte
	ReceiverLogin string
}

type UploaderService struct {
	Temp    rdb.TempGetter
	Storage storage.Storage
	waiters *waiter.WaiterService
}

type Uploader interface {
	StartFilesHTTPServer(port, storagePath string)
}

func NewUploaderService(temp rdb.TempGetter, storage storage.Storage, usedTempPaths map[string]struct{}) *UploaderService {
	return &UploaderService{
		Temp:    temp,
		Storage: storage,
		waiters: waiter.NewWaiterService(),
	}
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

	storagePath, err := s.Temp.GetStoragePath(ctx, r.URL.Path)
	if err != nil {
		slog.Error("Error get storage path", "error", err)
		w.Write([]byte("Error get storage path"))
		return
	}

	switch storagePath {
	case "save_avatar", "save_file":
		body, err := io.ReadAll(r.Body)
		if err != nil {
			s.notifyWaiter(r.URL.Path, "", err)

			slog.Error("Error read request body", "error", err)
			w.Write([]byte("Error read request body"))
			return
		}

		req := &FilesRequest{}
		if err = json.Unmarshal(body, &req); err != nil {
			s.notifyWaiter(r.URL.Path, "", err)

			w.Write([]byte("Error unmarshal request body"))
			return
		}

		saveAvatar := false
		if storagePath == "save_avatar" {
			saveAvatar = true
		}

		storagePath, err = s.Storage.SaveFile(req.File, req.ReceiverLogin, saveAvatar)
		if err != nil {
			s.notifyWaiter(r.URL.Path, "", err)

			slog.Error("Error save file", "error", err)
			w.Write([]byte("Error save file"))
			return
		}

		s.notifyWaiter(r.URL.Path, storagePath, nil)

		w.WriteHeader(http.StatusOK)
		return
	default:
		file, err := s.Storage.GetFile(storagePath)
		if err != nil {
			slog.Error("Cant read file", "error", err)
			w.Write([]byte("Error read file"))
			return
		}

		w.Write(file)
		return
	}
}

func (s *UploaderService) notifyWaiter(tempPath string, storagePath string, err error) {
	s.waiters.Mu.Lock()
	ch, ok := s.waiters.Waiters[tempPath]
	s.waiters.Mu.Unlock()

	if !ok {
		return
	}

	ch <- waiter.WaiterInfo{
		Path: storagePath,
		Err:  err,
	}
}
