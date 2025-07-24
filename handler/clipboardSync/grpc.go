package handler

import (
	"context"

	pb "github.com/YenXXXW/clipboradSyncServer/genproto/clipboardSync"
	"github.com/YenXXXW/clipboradSyncServer/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type ClipboardSypcGrpcHandler struct {
	clipboardSyncService types.ClipboardSyncService
	pb.UnimplementedClipSyncServiceServer
}

func NewGrpcClipboardSyncService(grpc *grpc.Server, clipboardSyncService types.ClipboardSyncService) {
	grpcHandler := &ClipboardSypcGrpcHandler{
		clipboardSyncService: clipboardSyncService,
	}

	pb.RegisterClipSyncServiceServer(grpc, grpcHandler)
}

func (h *ClipboardSypcGrpcHandler) SubscribeClipBoardContentUpdate(req *pb.SubscribeRequest, stream grpc.ServerStreamingServer[pb.ClipboardContent]) error {
	roomId := req.GetRoomId()
	deviceId := req.GetDeviceId()

	if err := h.clipboardSyncService.SubscribeClipBoardContentUpdate(deviceId, roomId, stream); err != nil {
		return status.Error(codes.Unavailable, err.Error())
	}

	return nil

}

func (h *ClipboardSypcGrpcHandler) SendClipBoardUpdate(ctx context.Context, req *pb.ClipboardUpdateRequest) (*emptypb.Empty, error) {

	err := h.clipboardSyncService.SendClipBoardUpdate(ctx, req.GetRoomId(), req.GetContent())
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &emptypb.Empty{}, nil
}
