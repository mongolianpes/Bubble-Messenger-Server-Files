package grpc

import (
	"context"
	"net"

	"files/internal/storage"
	pb "files/proto"

	"google.golang.org/grpc"
)

type GRPCService struct {
	pb.UnimplementedFilesServiceServer
	file storage.Storage
}

type GRPCServiceStarter interface {
	Start()
}

func NewSaveFilesService(storage storage.Storage) *GRPCService {
	return &GRPCService{
		file: storage,
	}
}

func (s *GRPCService) Save(ctx context.Context, req *pb.SaveFileRequest) (*pb.SaveFileResponse, error) {
	filepath, err := s.file.SaveFile(req.File)
	if err != nil {
		return nil, err
	}

	return &pb.SaveFileResponse{
		StoragePath: filepath,
	}, nil
}

func (s *GRPCService) DelFile(ctx context.Context, req *pb.DelFileRequest) (*pb.DelFileResponse, error) {
	if err := s.file.DelFile(req.StoragePath); err != nil {
		return nil, err
	}

	return &pb.DelFileResponse{}, nil
}

func (s *GRPCService) Start() {
	grpcServer := grpc.NewServer()

	lis, err := net.Listen("tcp", ":8086")
	if err != nil {
		panic(err)
	}

	pb.RegisterFilesServiceServer(grpcServer, NewSaveFilesService(s.file))
	if err := grpcServer.Serve(lis); err != nil {
		panic(err)
	}
}
