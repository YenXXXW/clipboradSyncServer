package handler

import (
	"context"
	"fmt"

	pb "github.com/YenXXXW/clipboradSyncServer/genproto/clipboardSync"
	"github.com/YenXXXW/clipboradSyncServer/shared"
	"github.com/YenXXXW/clipboradSyncServer/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type ClipboardSypcGrpcHandler struct {
	clipboardSyncService types.ClipboardSyncService
	roomService          types.RoomService
	pb.UnimplementedClipSyncServiceServer
}

func NewGrpcClipboardSyncService(grpc *grpc.Server, clipboardSyncService types.ClipboardSyncService, roomService types.RoomService) {
	grpcHandler := &ClipboardSypcGrpcHandler{
		clipboardSyncService: clipboardSyncService,
		roomService:          roomService,
	}

	pb.RegisterClipSyncServiceServer(grpc, grpcHandler)
}

func (h *ClipboardSypcGrpcHandler) SubscribeClipboardContentUpdate(req *pb.SubscribeRequest, stream grpc.ServerStreamingServer[pb.ClipboardUpdate]) error {
	roomId := req.GetRoomId()
	deviceId := req.GetDeviceId()

	if err := h.clipboardSyncService.SubscribeClipBoardContentUpdate(deviceId, roomId, &grpcStreamWrapper{stream: stream}); err != nil {
		return status.Error(codes.Unavailable, err.Error())
	}

	return nil

}

func (h *ClipboardSypcGrpcHandler) SendClipboardUpdate(ctx context.Context, req *pb.ClipboardUpdate) (*emptypb.Empty, error) {

	deviceID := req.GetDeviceId()
	client, ok := h.roomService.GetClient(deviceID)

	if !ok {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("client with the device_id: %s does not exist", deviceID))
	}

	content := shared.ClipboardContent{
		Text: req.GetContent().GetText(),
	}

	update := &shared.ClipboardUpdate{
		DeviceId: req.GetDeviceId(),
		Content:  content,
	}

	err := h.clipboardSyncService.SendClipBoardUpdate(ctx, client.RoomID, update)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (h *ClipboardSypcGrpcHandler) CreateRoom(ctx context.Context, req *pb.CreateRoomRequest) (*pb.CreateRoomResponse, error) {
	roomId := h.roomService.CreateRoom()

	res := &pb.CreateRoomResponse{
		RoomId: roomId,
	}

	return res, nil

}

func (h *ClipboardSypcGrpcHandler) LeaveRoom(ctx context.Context, req *pb.LeaveRoomRequest) (*emptypb.Empty, error) {

	h.roomService.DeleteClient(req.GetDeviceId())
	return &emptypb.Empty{}, nil

}
