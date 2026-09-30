package uploader

import (
	"files/internal/storage"
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
	fsFiles := http.FileServer(http.Dir(storagePath + storage.TimeToLiveFilesDirPath))
	fsAvatars := http.FileServer(http.Dir(storagePath + storage.ForeversDirPath))

	http.Handle("/users/", fsFiles)
	http.Handle("/avatars/", fsAvatars)

	err := http.ListenAndServe(port, nil)
	if err != nil {
		slog.Error("Error start HTTP uploader", "error", err)
	}
}
