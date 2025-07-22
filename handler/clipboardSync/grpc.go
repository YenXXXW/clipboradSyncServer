package handler

import (
	"context"

	pb "github.com/YenXXXW/clipboradSyncServer/genproto/clipboardSync"
	"github.com/YenXXXW/clipboradSyncServer/types"
	"google.golang.org/grpc"

	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type ClipboardGrpcHandler struct {
	clipboardSyncService types.ClipboardSyncService
	pb.UnimplementedClipSyncServiceServer
}

func (h *ClipboardGrpcHandler) SubscribeClipBoardContentUpdate(req *pb.SubscribeRequest, stream grpc.ServerStreamingServer[pb.ClipboardContent]) error {
	return nil

}

func (h *ClipboardGrpcHandler) SendClipBoardUpdate(ctx context.Context, data *pb.ClipboardContent) (*emptypb.Empty, error) {

	return nil, nil
}
