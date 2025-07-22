package service

import (
	"context"

	pb "github.com/YenXXXW/clipboradSyncServer/genproto/clipboardSync"
	"github.com/YenXXXW/clipboradSyncServer/types"

	"google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type ClipboardSyncService struct {
	RoomManager   *RoomManager
	ClientManager *ClientManager
}

func NewClipboardSyncService() *ClipboardSyncService {
	return &ClipboardSyncService{}
}

func (s *ClipboardSyncService) SendClipBoardUpdate(ctx context.Context, req *pb.ClipboardUpdateRequest) (*emptypb.Empty, error) {
	roomID := req.GetRoomId()
	content := req.GetContent()

	s.RoomManager.BroadcastToRoom(roomID, content)

	return nil, nil
}

func (s *ClipboardSyncService) SubscribeClipBoardContentUpdate(req *pb.SubscribeRequest, grpc grpc.ServerStreamingServer[pb.ClipboardContent]) error {
	deviceId := req.GetDeviceId()
	roomId := req.GetRoomId()

	var client *types.Client
	exitingClient, ok := s.ClientManager.Clients[deviceId]
	if ok {
		client = exitingClient
		if client.RoomID == roomId {
			s.RoomManager.RemoveFromRoom(client.RoomID, client)
			client.RoomID = roomId
		}
	} else {
		client = s.ClientManager.CreateNewClient(roomId, grpc)
	}

	s.RoomManager.JoinRoom(roomId, client)

	for {
		select {
		case msg := <-client.Send:
			if err := client.Conn.SendMsg(msg); err != nil {
				s.RoomManager.RemoveFromRoom(client.RoomID, client)
				s.ClientManager.RemoveClient(client.ID)
				return err
			}

		case <-client.Conn.Context().Done():
			s.RoomManager.RemoveFromRoom(client.RoomID, client)
			s.ClientManager.RemoveClient(client.ID)
			return nil

		}
	}
}
