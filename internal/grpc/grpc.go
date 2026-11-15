package grpc

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"log/slog"
	"net"

	"files/internal/rdb"
	"files/internal/storage"
	pb "files/proto"

	"google.golang.org/grpc"
)

type GRPCService struct {
	pb.UnimplementedFilesServiceServer
	file storage.Storage
	temp rdb.TempSetter
}

type GRPCServiceStarter interface {
	Start()
}

func NewSaveFilesService(storage storage.Storage) *GRPCService {
	return &GRPCService{
		file: storage,
	}
}

func (s *GRPCService) Del(ctx context.Context, req *pb.DelFileRequest) (*pb.DelFileResponse, error) {
	if err := s.file.DelFile(req.StoragePath); err != nil {
		slog.ErrorContext(ctx, "Cant del file", "storagePath", req.StoragePath, "error", err)
		return nil, err
	}

	return &pb.DelFileResponse{}, nil
}

func (s *GRPCService) GenerateTempPath(ctx context.Context, req *pb.GenerateTempPathReqeust) (*pb.GenerateTempPathResponse, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}

	tempPath := base64.RawURLEncoding.EncodeToString(b)

	if err := s.temp.SetTempPath(ctx, tempPath, req.Info); err != nil {
		slog.ErrorContext(ctx, "Cant save temp path", "info", req.Info, "error", err)
		return nil, err
	}

	return &pb.GenerateTempPathResponse{TempPath: tempPath}, nil
}

func (s *GRPCService) Start() {
	grpcServer := grpc.NewServer()

	lis, err := net.Listen("tcp", ":8086")
	if err != nil {
		slog.Error("Error setup net listener for gRPC service", "error", err)
	}

	pb.RegisterFilesServiceServer(grpcServer, NewSaveFilesService(s.file))
	if err := grpcServer.Serve(lis); err != nil {
		slog.Error("Error serve gRPC service", "error", err)
	}
}
