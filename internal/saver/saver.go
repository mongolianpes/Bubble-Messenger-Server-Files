package saver

import (
	"context"

	"files/internal/storage"
	pb "files/proto"
)

type SaveFilesService struct {
	pb.UnimplementedFilesServiceServer
	storage storage.Saver
}

type Saver interface {
	File(ctx context.Context, req *pb.SaveFileRequest) (*pb.SaveFileResponse, error)
}

func NewSaveFilesService(storage storage.Saver) *SaveFilesService {
	return &SaveFilesService{
		storage: storage,
	}
}

func (s *SaveFilesService) File(ctx context.Context, req *pb.SaveFileRequest) (*pb.SaveFileResponse, error) {
	filepath, err := s.storage.Save(req.File)
	if err != nil {
		return nil, err
	}

	return &pb.SaveFileResponse{
		StoragePath: filepath,
	}, nil
}
