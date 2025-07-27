package handler

import (
	pb "github.com/YenXXXW/clipboradSyncServer/genproto/clipboardSync"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type grpcStreamWrapper struct {
	stream grpc.ServerStreamingServer[pb.ClipboardContent]
}

func (g *grpcStreamWrapper) Send(content *pb.ClipboardContent) error {
	return g.stream.Send(content)
}

func (g *grpcStreamWrapper) Context() context.Context {
	return g.stream.Context()
}
