package downloader

import (
	"context"

	"files/internal/storage"
	pb "files/proto"
)

type DonwloaderFilesService struct {
	pb.UnimplementedFilesServiceServer
	storage storage.Saver
}

type Downloader interface {
	DownloadFile(ctx context.Context, req *pb.DownloadFileRequest) (*pb.DownloadFileResponse, error)
}

func NewDonwloaderFilesService(storage storage.Saver) *DonwloaderFilesService {
	return &DonwloaderFilesService{
		storage: storage,
	}
}

func (s *DonwloaderFilesService) DownloadFile(ctx context.Context, req *pb.DownloadFileRequest) (*pb.DownloadFileResponse, error) {
	filepath, err := s.storage.Save(req.File)
	if err != nil {
		return nil, err
	}

	return &pb.DownloadFileResponse{
		StoragePath: filepath,
	}, nil
}
