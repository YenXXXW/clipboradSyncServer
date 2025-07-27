package service

import (
	"context"
	"log"

	pb "github.com/YenXXXW/clipboradSyncServer/genproto/clipboardSync"
	"github.com/YenXXXW/clipboradSyncServer/shared"
	"github.com/YenXXXW/clipboradSyncServer/types"
)

type ClipboardSyncService struct {
	RoomService types.RoomService
}

func NewClipboardSyncService(roomService types.RoomService) *ClipboardSyncService {
	return &ClipboardSyncService{
		RoomService: roomService,
	}
}

func (s *ClipboardSyncService) SendClipBoardUpdate(ctx context.Context, roomID string, content *pb.ClipboardContent) error {

	if err := s.RoomService.BroadcastToRoom(roomID, content); err != nil {
		return err
	}
	return nil

}

func (s *ClipboardSyncService) SubscribeClipBoardContentUpdate(deviceId, roomId string, stream shared.StreamWriter) error {

	var client *types.Client
	existingClient, ok := s.RoomService.GetClient(deviceId)
	if ok {
		client = existingClient
		if client.RoomID != roomId {
			s.RoomService.RemoveFromRoom(deviceId, client.RoomID)
			client.RoomID = roomId
		}
	} else {

		client = s.RoomService.CreateClient(deviceId, roomId, stream)
	}

	s.RoomService.JoinRoom(roomId, client)

	log.Printf("device: %s joined to the room: %s", deviceId, roomId)

	for {
		select {
		case msg := <-client.Send:
			if err := client.Conn.Send(msg); err != nil {
				s.RoomService.RemoveFromRoom(deviceId, client.RoomID)
				s.RoomService.DeleteClient(client.ID)
				return err
			}

		case <-client.Conn.Context().Done():
			s.RoomService.RemoveFromRoom(deviceId, client.RoomID)
			s.RoomService.DeleteClient(client.ID)
			return nil

		}
	}
}
