package service

import (
	"log/slog"
	"time"

	"files/internal/grpc"
	"files/internal/rdb"
	"files/internal/storage"
	"files/internal/uploader"
	"files/internal/waiter"
)

const filesDirPath = "storage/"

type FilesService struct {
	uploader    uploader.Uploader
	grpcService grpc.GRPCServiceStarter
	cleaner     storage.Cleaner
}

func NewFilesService() *FilesService {
	filesStorage := storage.NewStorage(filesDirPath)
	temp, err := rdb.NewRedisStorage()
	if err != nil {
		slog.Error("Failed to start Redis storage", "error", err)
	}
	waiters := waiter.NewWaiterService()
	grpcService := grpc.NewSaveFilesService(filesStorage, temp, waiters)
	uploaderService := uploader.NewUploaderService(temp, filesStorage, waiters)
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
	s.uploader.StartFilesHTTPServer(":8080", filesDirPath)
}

func (s *FilesService) StartFilesCleaner() {
	for {
		s.cleaner.RemoveOld()

		time.Sleep(time.Hour * 25)
	}
}
