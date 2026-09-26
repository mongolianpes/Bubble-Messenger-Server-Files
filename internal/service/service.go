package service

import (
	"net"
	"time"

	"files/internal/downloader"
	"files/internal/storage"
	"files/internal/uploader"
	pb "files/proto"

	"google.golang.org/grpc"
)

type FilesService struct {
	uploader   uploader.Uploader
	downloader downloader.Downloader
	storage    storage.Cleaner
}

func NewFilesService() *FilesService {
	return &FilesService{}
}

func (s *FilesService) StartDownloader() {
	filesStorage := storage.NewStorage()
	grpcServer := grpc.NewServer()

	lis, err := net.Listen("tcp", ":8086")
	if err != nil {
		panic(err)
	}

	pb.RegisterFilesServiceServer(grpcServer, downloader.NewDonwloaderFilesService(filesStorage))
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}

func (s *FilesService) StartUploader() {
	uploaderService := uploader.NewUploaderService()
	uploaderService.StartUploaderHHTPservice(":8080")
}

func (s *FilesService) StartFilesCleaner() {
	for {
		s.storage.RemoveOld()

		time.Sleep(time.Hour * 25)
	}
}
