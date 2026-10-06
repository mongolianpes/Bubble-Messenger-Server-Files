package uploader

import (
	"log/slog"
	"net/http"
)

type UploaderService struct{}

type Uploader interface {
	StartUploaderHHTPservice(port, storagePath string)
}

func NewUploaderService() *UploaderService {
	return &UploaderService{}
}

func (s *UploaderService) StartUploaderHHTPservice(port, storagePath string) {
	fs := http.FileServer(http.Dir(storagePath))

	http.Handle("/storage/", http.StripPrefix("/storage/", fs))

	err := http.ListenAndServe(port, nil)
	if err != nil {
		slog.Error("Error start HTTP uploader", "error", err)
	}
}
