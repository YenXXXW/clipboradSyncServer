package handler

import (
	pb "github.com/YenXXXW/clipboradSyncServer/genproto/clipboardSync"
	"github.com/YenXXXW/clipboradSyncServer/types"
	"google.golang.org/grpc"
)

type ClipboardGrpcHandler struct {
	clipboardSyncService types.ClipboardSyncService
	pb.UnimplementedClipSyncServiceServer
}

func (h *ClipboardGrpcHandler) SubscribeClipBoardContentUpdate(req *pb.SubscribeRequest, stream grpc.ServerStreamingServer[pb.ClipBoardContent]) error {
	return nil

}
