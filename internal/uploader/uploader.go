package uploader

import (
	"files/internal/storage"
	"log"
	"net/http"
)

type UploaderService struct{}

type Uploader interface {
	StartUploaderHHTPservice(port string)
}

func NewUploaderService() *UploaderService {
	return &UploaderService{}
}

func (s *UploaderService) StartUploaderHHTPservice(port string) {
	fs := http.FileServer(http.Dir(storage.FilesDirPath))

	http.Handle("/", fs)

	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
