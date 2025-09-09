package service

import (
	"context"
	"errors"
	"fmt"
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

	client, clientExists := s.RoomService.GetClient(deviceId)

	_, roomExists := s.RoomService.GetRoom(roomId)
	fmt.Println()

	if !roomExists {

		validateJoin := &shared.ValidateJoin{}
		roomCheckError := shared.Validate{
			Success: false,
			Message: "Room does not exist",
		}

		validateJoin.ValidateRoom = roomCheckError

		if err := stream.Send(&shared.UpdateEvent{

			ValidateJoin: validateJoin,
		}); err != nil {

			log.Println("Error sending Validate Join Message", err)
			return err
		}

		return errors.New("room with given ID does not exist")

	}

	if clientExists && client.RoomID != "" {

		validateJoin := &shared.ValidateJoin{}
		checkClientError := shared.Validate{
			Success: false,
			Message: "Client is in a room",
		}
		validateJoin.CheckClient = checkClientError

		if err := stream.Send(&shared.UpdateEvent{
			ValidateJoin: validateJoin,
		}); err != nil {
			log.Println("Error sending Validate Message when roomId is not empty string", err)
			return err
		}

		return errors.New("client is still in a room")
	}

	if !clientExists {

		client = s.RoomService.CreateClient(deviceId, roomId)
	}

	UpdateEvent := &shared.UpdateEvent{
		ValidateJoin: &shared.ValidateJoin{
			ValidateRoom: shared.Validate{
				Success: true,
				Message: "Room Validate Successful",
			},
			CheckClient: shared.Validate{
				Success: true,
				Message: "Client Validate Successful",
			},
		},
	}

	if err := stream.Send(UpdateEvent); err != nil {

		log.Println("Error sending message on stream", err)
		return err
	}

	//if !roomExists || (clientExists && client.RoomID != "") {
	//
	//log.Println("entered this loope")
	//validateJoin := &shared.ValidateJoin{}
	//if !roomExists {
	//roomCheckError := shared.Validate{
	//Success: false,
	//Message: "Room does not exist",
	//}
	//
	//validateJoin.ValidateRoom = roomCheckError
	//
	//}
	//
	//if clientExists && client.RoomID != "" {
	//checkClientError := shared.Validate{
	//Success: false,
	//Message: "Client is in a room",
	//}
	//validateJoin.CheckClient = checkClientError
	//}
	//
	//log.Println("works fine ")
	//
	//log.Println("working till now")
	//
	//if err := client.Conn.Send(&shared.UpdateEvent{
	//ValidateJoin: validateJoin,
	//}); err != nil {
	//log.Println("Error sending message on stream")
	//return err
	//}
	//
	//return errors.New("Room does not exist or the client has joined a room")
	//}
	//
	//if !clientExists {
	//client = s.RoomService.CreateClient(deviceId, roomId, stream)
	//}

	s.RoomService.JoinRoom(roomId, client)

	validateJoin := &shared.ValidateJoin{
		CheckClient: shared.Validate{
			Success: true,
		},
		ValidateRoom: shared.Validate{
			Success: true,
		},
	}

	if err := stream.Send(&shared.UpdateEvent{
		ValidateJoin: validateJoin,
	}); err != nil {
		log.Println("Error sending message on stream", err)
		return err
	}

	log.Printf("device: %s joined to the room: %s", deviceId, roomId)

	for {
		select {
		case msg := <-client.Send:
			update := &shared.UpdateEvent{
				ClipboardUpdate: msg,
			}
			if err := stream.Send(update); err != nil {
				s.RoomService.DeleteClient(client.DeviceID)
				return err
			}
		case <-stream.Context().Done():
			//when client calls context.done delete the client struct
			s.RoomService.DeleteClient(client.DeviceID)
			return nil
		case <-client.Done:
			log.Printf("device: %s disconnected from room: %s", deviceId, roomId)
			return nil
		}
	}
}
