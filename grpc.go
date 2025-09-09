package main

import (
	"log"
	"net"

	handler "github.com/YenXXXW/clipboradSyncServer/handler/clipboardSync"
	"github.com/YenXXXW/clipboradSyncServer/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type gRPCServer struct {
	addr string
}

func NewgRPCServer(addr string) *gRPCServer {
	return &gRPCServer{
		addr: addr,
	}
}

func (s *gRPCServer) Run() error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	roomService := service.NewRoomService()
	clipboardSyncService := service.NewClipboardSyncService(roomService)
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(handler.UnaryInterceptor), grpc.StreamInterceptor(handler.StreamInterceptor))
	handler.NewGrpcClipboardSyncService(grpcServer, clipboardSyncService, roomService)

	reflection.Register(grpcServer)

	log.Println("Starting the gRPC server on", s.addr)
	return grpcServer.Serve(lis)
}
