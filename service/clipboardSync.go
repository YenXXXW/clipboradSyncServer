package service

import (
	"context"
	"log"

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

func (s *ClipboardSyncService) SendClipBoardUpdate(ctx context.Context, roomID string, content *shared.ClipboardUpdate) error {

	if err := s.RoomService.BroadcastToRoom(roomID, content); err != nil {
		return err
	}
	return nil

}

func (s *ClipboardSyncService) SubscribeClipBoardContentUpdate(deviceId, roomId string, stream shared.StreamWriter) error {

	s.RoomService.DeleteClient(deviceId)

	client := s.RoomService.CreateClient(deviceId, roomId, stream)
	s.RoomService.JoinRoom(roomId, client)

	log.Printf("device: %s joined to the room: %s", deviceId, roomId)

	for {
		select {
		case msg := <-client.Send:
			if err := client.Conn.Send(msg); err != nil {
				s.RoomService.DeleteClient(client.DeviceID)
				return err
			}
		case <-client.Conn.Context().Done():
			s.RoomService.DeleteClient(client.DeviceID)
			return nil
		case <-client.Done:
			log.Printf("device: %s disconnected from room: %s", deviceId, roomId)
			return nil
		}
	}
}
