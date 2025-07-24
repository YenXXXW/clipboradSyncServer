package types

import (
	"context"
	pb "github.com/YenXXXW/clipboradSyncServer/genproto/clipboardSync"
	"google.golang.org/grpc"
)

type ClipboardSyncService interface {
	SubscribeClipBoardContentUpdate(string, string, grpc.ServerStreamingServer[pb.ClipboardContent]) error
	SendClipBoardUpdate(context.Context, string, *pb.ClipboardContent) error
}
