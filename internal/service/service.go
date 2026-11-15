package service

import (
	"time"

	"files/internal/grpc"
	"files/internal/rdb"
	"files/internal/storage"
	"files/internal/uploader"
)

const filesDirPath = "storage/"

type FilesService struct {
	uploader    uploader.Uploader
	grpcService grpc.GRPCServiceStarter
	cleaner     storage.Cleaner
}

func NewFilesService() *FilesService {
	usedTempPaths := make(map[string]struct{})
	filesStorage := storage.NewStorage(filesDirPath)
	grpcService := grpc.NewSaveFilesService(filesStorage, usedTempPaths)
	temp, _ := rdb.NewRedisStorage()
	uploaderService := uploader.NewUploaderService(temp, filesStorage, usedTempPaths)
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
