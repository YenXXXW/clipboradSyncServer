package service

import (
	"context"

	pb "github.com/YenXXXW/clipboradSyncServer/genproto/clipboardSync"
	"github.com/YenXXXW/clipboradSyncServer/types"

	"google.golang.org/grpc"
)

type ClipboardSyncService struct {
	RoomManager   *RoomManager
	ClientManager *ClientManager
}

func NewClipboardSyncService() *ClipboardSyncService {
	newRoomManager := RoomManager{
		Rooms: make(map[string]*types.Room),
	}

	newClientManager := ClientManager{
		Clients: make(map[string]*types.Client),
	}
	return &ClipboardSyncService{
		RoomManager:   &newRoomManager,
		ClientManager: &newClientManager,
	}
}

func (s *ClipboardSyncService) SendClipBoardUpdate(ctx context.Context, roomID string, content *pb.ClipboardContent) error {

	if err := s.RoomManager.BroadcastToRoom(roomID, content); err != nil {
		return err
	}
	return nil

}

func (s *ClipboardSyncService) SubscribeClipBoardContentUpdate(deviceId, roomId string, grpc grpc.ServerStreamingServer[pb.ClipboardContent]) error {

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
