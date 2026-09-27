package service

import (
	"time"

	"files/internal/grpc"
	"files/internal/storage"
	"files/internal/uploader"
)

type FilesService struct {
	uploader    uploader.Uploader
	grpcService grpc.GRPCServiceStarter
	cleaner     storage.Cleaner
}

func NewFilesService() *FilesService {
	filesStorage := storage.NewStorage()
	grpcService := grpc.NewSaveFilesService(filesStorage)
	uploaderService := uploader.NewUploaderService()
	return &FilesService{
		uploader:    uploaderService,
		grpcService: grpcService,
		cleaner:     filesStorage,
	}
}

func (s *FilesService) StartGRPCService() {
	s.grpcService.Start()
}

func (s *FilesService) StartUploader() {
	s.uploader.StartUploaderHHTPservice(":8080")
}

func (s *FilesService) StartFilesCleaner() {
	for {
		s.cleaner.RemoveOld()

		time.Sleep(time.Hour * 25)
	}
}
