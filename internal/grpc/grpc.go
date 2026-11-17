package grpc

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log/slog"
	"net"

	"files/internal/rdb"
	"files/internal/storage"
	"files/internal/waiter"
	pb "files/proto"

	"google.golang.org/grpc"
)

const (
	saveAvatar = "save_avatar"
	saveFile   = "save_file"
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

func NewSaveFilesService(storage storage.Storage, temp rdb.TempSetter, waiters *waiter.WaiterService) *GRPCService {
	return &GRPCService{
		file:    storage,
		temp:    temp,
		waiters: waiters,
	}
}

func (s *GRPCService) Del(ctx context.Context, req *pb.DelFileRequest) (*pb.DelFileResponse, error) {
	if err := s.file.DelFile(req.StoragePath); err != nil {
		slog.ErrorContext(ctx, "Cant del file", "storagePath", req.StoragePath, "error", err)
		return nil, err
	}

	return &pb.DelFileResponse{}, nil
}

func (s *GRPCService) TempPath(req *pb.TempPathReqeust, stream pb.FilesService_TempPathServer) error {
	ctx := stream.Context()

	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return err
	}

	tempPath := "/" + base64.RawURLEncoding.EncodeToString(b)

	switch req.Info {
	case saveAvatar, saveFile:
	default:
		if !s.file.Exist(req.Info) {
			return fmt.Errorf("File %s not exist", req.Info)
		}
	}

	if err := s.temp.SetStoragePath(ctx, tempPath, req.Info); err != nil {
		slog.ErrorContext(ctx, "Cant save temp path", "info", req.Info, "error", err)
		return err
	}

	switch req.Info {
	case "save_avatar", "save_file":
	default:
		if err := stream.Send(&pb.TempPathResponse{TempPath: tempPath}); err != nil {
			return err
		}

		return nil
	}

	ch := make(chan waiter.WaiterInfo, 1)

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
	case status := <-ch:
		if status.Err != nil {
			return status.Err
		}
		if err := stream.Send(&pb.TempPathResponse{StoragePath: status.Path}); err != nil {
			return err
		}
	case <-ctx.Done():
		return fmt.Errorf("Ошибка использования канала")
	}

	return nil
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
