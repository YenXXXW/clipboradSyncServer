package service

import (
	"errors"
	"fmt"
	"sync"

	pb "github.com/YenXXXW/clipboradSyncServer/genproto/clipboardSync"
	"github.com/YenXXXW/clipboradSyncServer/shared"
	"github.com/YenXXXW/clipboradSyncServer/types"
	"github.com/google/uuid"
)

type RoomManager struct {
	Rooms map[string]*types.Room
	Mutex sync.Mutex
}

type ClientManager struct {
	Clients map[string]*types.Client
	Mutex   sync.Mutex
}

type RoomService struct {
	clientManager *ClientManager
	roomManager   *RoomManager
}

func NewRoomService() *RoomService {

	roomManager := &RoomManager{
		Rooms: make(map[string]*types.Room),
	}

	clientManager := &ClientManager{
		Clients: make(map[string]*types.Client),
	}

	return &RoomService{
		roomManager:   roomManager,
		clientManager: clientManager,
	}
}

func NewRoom() *types.Room {
	return &types.Room{
		Clients:          make(map[string]*types.Client),
		ClipBoardContent: "",
	}
}

func NewClient(roomID string, conn shared.StreamWriter) *types.Client {
	return &types.Client{
		ID:     uuid.NewString(),
		RoomID: roomID,
		Conn:   conn,
		Send:   make(chan *pb.ClipboardContent, 20),
	}
}

func (cm *ClientManager) CreateNewClient(deviceID, roomID string, conn shared.StreamWriter) *types.Client {
	cm.Mutex.Lock()
	defer cm.Mutex.Unlock()
	newClient := NewClient(roomID, conn)
	cm.Clients[deviceID] = newClient
	return newClient
}

func (cm *ClientManager) GetClient(deviceID string) (*types.Client, bool) {
	cm.Mutex.Lock()
	defer cm.Mutex.Unlock()
	client, ok := cm.Clients[deviceID]
	return client, ok
}

func (cm *ClientManager) RemoveClient(deviceID string) {
	cm.Mutex.Lock()
	defer cm.Mutex.Unlock()
	delete(cm.Clients, deviceID)
}

func (rm *RoomManager) CreateRoom() string {
	rm.Mutex.Lock()
	defer rm.Mutex.Unlock()
	id := uuid.New().String()
	newRoom := NewRoom()
	rm.Rooms[id] = newRoom
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

func (rm *RoomManager) BroadcastToRoom(roomID string, clipboardData *pb.ClipboardContent) error {
	rm.Mutex.Lock()
	defer rm.Mutex.Unlock()
	room, ok := rm.Rooms[roomID]
	if !ok {
		return fmt.Errorf("room %s not found", roomID)
	}
	for _, client := range room.Clients {
		select {
		case client.Send <- clipboardData:

		default:
			rm.RemoveFromRoom(roomID, client)
		}
	}

	return nil

}

func (s *RoomService) CreateRoom() string {
	return s.roomManager.CreateRoom()
}

func (s *RoomService) JoinRoom(roomID string, client *types.Client) error {
	return s.roomManager.JoinRoom(roomID, client)
}

func (s *RoomService) RemoveFromRoom(deviceID, roomID string) error {
	client, ok := s.clientManager.GetClient(deviceID)
	if !ok {
		return errors.New("client with the give device_id does not exist")
	}
	return s.roomManager.RemoveFromRoom(roomID, client)
}

func (s *RoomService) CreateClient(deviceID, roomID string, conn shared.StreamWriter) *types.Client {

	return s.clientManager.CreateNewClient(deviceID, roomID, conn)
}

func (s *RoomService) GetClient(deviceID string) (*types.Client, bool) {
	return s.clientManager.GetClient(deviceID)
}

func (s *RoomService) DeleteClient(deviceID string) {
	s.clientManager.RemoveClient(deviceID)
}

func (s *RoomService) BroadcastToRoom(roomID string, clipboardData *pb.ClipboardContent) error {
	return s.roomManager.BroadcastToRoom(roomID, clipboardData)
}
