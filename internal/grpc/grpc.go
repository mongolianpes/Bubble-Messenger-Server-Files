package grpc

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"log/slog"
	"net"

	"files/internal/rdb"
	"files/internal/storage"
	"files/internal/waiter"
	pb "files/proto"

	"google.golang.org/grpc"
)

type GRPCService struct {
	pb.UnimplementedFilesServiceServer
	file    storage.Storage
	temp    rdb.TempSetter
	waiters *waiter.WaiterService
}

type GRPCServiceStarter interface {
	Start()
}

func NewSaveFilesService(storage storage.Storage, usedTempPaths map[string]struct{}) *GRPCService {
	return &GRPCService{
		file:    storage,
		waiters: waiter.NewWaiterService(),
	}
}

func (s *GRPCService) Del(ctx context.Context, req *pb.DelFileRequest) (*pb.DelFileResponse, error) {
	if err := s.file.DelFile(req.StoragePath); err != nil {
		slog.ErrorContext(ctx, "Cant del file", "storagePath", req.StoragePath, "error", err)
		return nil, err
	}

	return &pb.DelFileResponse{}, nil
}

func (s *GRPCService) GenerateTempPath(req *pb.TempPathReqeust, stream pb.FilesService_TempPathServer) error {
	ctx := stream.Context()

	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return err
	}

	tempPath := base64.RawURLEncoding.EncodeToString(b)

	if err := s.temp.SetTempPath(ctx, tempPath, req.Info); err != nil {
		slog.ErrorContext(ctx, "Cant save temp path", "info", req.Info, "error", err)
		return err
	}

	ch := make(chan error, 1)

	s.waiters.Mu.Lock()
	s.waiters.Waiters[tempPath] = ch
	s.waiters.Mu.Unlock()

	defer func() {
		s.waiters.Mu.Lock()
		delete(s.waiters.Waiters, tempPath)
		s.waiters.Mu.Unlock()
	}()

	if err := stream.Send(&pb.TempPathResponse{TempPath: tempPath}); err != nil {
		return err
	}

	select {
	case err := <-ch:
		if err != nil {
			return stream.Send(&pb.TempPathResponse{Success: false})
		}
		return stream.Send(&pb.TempPathResponse{Success: true})
	case <-ctx.Done():
		return stream.Send(&pb.TempPathResponse{Success: false})
	}
}

func (s *GRPCService) Start() {
	grpcServer := grpc.NewServer()

	lis, err := net.Listen("tcp", ":8086")
	if err != nil {
		slog.Error("Error setup net listener for gRPC service", "error", err)
	}

	pb.RegisterFilesServiceServer(grpcServer, s)
	if err := grpcServer.Serve(lis); err != nil {
		slog.Error("Error serve gRPC service", "error", err)
	}
}
