package service

import (
	"fmt"
	"sync"

	pb "github.com/YenXXXW/clipboradSyncServer/genproto/clipboardSync"
	"github.com/YenXXXW/clipboradSyncServer/types"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

func NewRoom() *types.Room {
	return &types.Room{
		ClipBoardContent: "",
	}
}

func NewClient(roomID string, conn grpc.ServerStreamingServer[pb.ClipboardContent]) *types.Client {
	return &types.Client{
		ID:     uuid.NewString(),
		RoomID: roomID,
		Conn:   conn,
		Send:   make(chan *pb.ClipboardContent),
	}
}

type RoomManager struct {
	Rooms map[string]*types.Room
	Mutex sync.Mutex
}

type ClientManager struct {
	Clients map[string]*types.Client
	Mutex   sync.Mutex
}

func (cm *ClientManager) CreateNewClient(roomID string, conn grpc.ServerStreamingServer[pb.ClipboardContent]) *types.Client {
	cm.Mutex.Lock()
	defer cm.Mutex.Unlock()
	newClient := NewClient(roomID, conn)
	*cm.Clients[newClient.ID] = *newClient
	return newClient
}

func (cm *ClientManager) RemoveClient(clientID string) {
	cm.Mutex.Lock()
	defer cm.Mutex.Unlock()
	delete(cm.Clients, clientID)
}

func (rm *RoomManager) CreateRoom() string {
	rm.Mutex.Lock()
	defer rm.Mutex.Unlock()
	id := uuid.New().String()
	newRoom := NewRoom()
	*rm.Rooms[id] = *newRoom
	return id
}

func (rm *RoomManager) JoinRoom(roomID string, client *types.Client) error {
	rm.Mutex.Lock()
	defer rm.Mutex.Unlock()

	room, ok := rm.Rooms[roomID]
	if !ok {
		return fmt.Errorf("room with ID %s does not exist", roomID)
	}

	room.Clients[client.ID] = client
	return nil
}

func (rm *RoomManager) RemoveFromRoom(roomID string, client *types.Client) error {
	rm.Mutex.Lock()
	defer rm.Mutex.Unlock()

	room, ok := rm.Rooms[roomID]
	if !ok {
		return fmt.Errorf("room with ID %s does not exist", roomID)
	}

	delete(room.Clients, client.ID)
	return nil
}

func (rm *RoomManager) BroadcastToRoom(roomID string, clipboardData *pb.ClipboardContent) {
	rm.Mutex.Lock()
	defer rm.Mutex.Unlock()
	if room, ok := rm.Rooms[roomID]; ok {
		for _, client := range room.Clients {
			select {
			case client.Send <- clipboardData:

			default:
				rm.RemoveFromRoom(roomID, client)
			}
		}

	}
}
