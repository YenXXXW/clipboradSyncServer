package handler

import (
	pb "github.com/YenXXXW/clipboradSyncServer/genproto/clipboardSync"
	"github.com/YenXXXW/clipboradSyncServer/shared"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type grpcStreamWrapper struct {
	stream grpc.ServerStreamingServer[pb.ClipboardUpdate]
}

func (g *grpcStreamWrapper) Send(update *shared.ClipboardUpdate) error {
	content := &pb.ClipboardContent{
		Text: update.Content.Text,
	}
	grpcUpdate := &pb.ClipboardUpdate{
		DeviceId: update.DeviceId,
		Content:  content,
	}
	return g.stream.Send(grpcUpdate)
}

func (g *grpcStreamWrapper) Context() context.Context {
	return g.stream.Context()
}
