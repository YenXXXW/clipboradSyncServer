package service

import (
	"errors"
	"fmt"
	"sync"

	"github.com/YenXXXW/clipboradSyncServer/shared"
	"github.com/YenXXXW/clipboradSyncServer/types"
	"github.com/google/uuid"
)

type roomManager struct {
	rooms map[string]*types.Room
	mutex sync.Mutex
}

type clientManager struct {
	clients map[string]*types.Client
	mutex   sync.Mutex
}

type RoomService struct {
	clientManager *clientManager
	roomManager   *roomManager
}

func NewRoomService() *RoomService {
	return &RoomService{
		roomManager: &roomManager{
			rooms: make(map[string]*types.Room),
		},
		clientManager: &clientManager{
			clients: make(map[string]*types.Client),
		},
	}
}

// --- Factory methods ---

func newRoom() *types.Room {
	return &types.Room{
		Clients:          make(map[string]*types.Client),
		ClipBoardContent: "",
	}
}

func newClient(deviceID, roomID string, conn shared.StreamWriter) *types.Client {
	return &types.Client{
		ID:       uuid.NewString(),
		DeviceID: deviceID,
		RoomID:   roomID,
		Conn:     conn,
		Send:     make(chan *shared.ClipboardUpdate, 20),
		Done:     make(chan struct{}),
	}
}

// --- RoomService methods  ---

// CreateRoom creates a new room and returns its ID.
func (s *RoomService) CreateRoom() string {
	s.roomManager.mutex.Lock()
	defer s.roomManager.mutex.Unlock()

	id := uuid.New().String()
	s.roomManager.rooms[id] = newRoom()
	return id
}

// JoinRoom adds a client to a specified room.
func (s *RoomService) JoinRoom(roomID string, client *types.Client) error {
	s.roomManager.mutex.Lock()
	defer s.roomManager.mutex.Unlock()

	room, ok := s.roomManager.rooms[roomID]
	if !ok {
		return fmt.Errorf("room with ID %s does not exist", roomID)
	}
	room.Clients[client.ID] = client
	return nil
}

// RemoveFromRoom removes a client from a room.
func (s *RoomService) RemoveFromRoom(deviceID, roomID string) error {
	s.roomManager.mutex.Lock()
	defer s.roomManager.mutex.Unlock()
	s.clientManager.mutex.Lock()
	defer s.clientManager.mutex.Unlock()

	client, ok := s.clientManager.clients[deviceID]
	if !ok {
		return errors.New("client with the given device_id does not exist")
	}

	room, ok := s.roomManager.rooms[roomID]
	if ok {
		delete(room.Clients, client.ID)
	}
	return nil
}

func (s *RoomService) GetRoom(roomID string) (*types.Room, bool) {
	s.roomManager.mutex.Lock()
	defer s.roomManager.mutex.Unlock()

	room, ok := s.roomManager.rooms[roomID]
	return room, ok
}

// CreateClient creates a new client and adds it to the client manager.
func (s *RoomService) CreateClient(deviceID, roomID string, conn shared.StreamWriter) *types.Client {
	s.clientManager.mutex.Lock()
	defer s.clientManager.mutex.Unlock()

	newClient := newClient(deviceID, roomID, conn)
	s.clientManager.clients[deviceID] = newClient
	return newClient
}

// GetClient retrieves a client by their device ID.
func (s *RoomService) GetClient(deviceID string) (*types.Client, bool) {
	s.clientManager.mutex.Lock()
	defer s.clientManager.mutex.Unlock()

	client, ok := s.clientManager.clients[deviceID]
	return client, ok
}

// DeleteClient removes a client from the system entirely. It removes the client
// from its room, terminates its subscription goroutine, and deletes it from the manager.
func (s *RoomService) DeleteClient(deviceID string) {
	s.clientManager.mutex.Lock()
	defer s.clientManager.mutex.Unlock()
	s.roomManager.mutex.Lock()
	defer s.roomManager.mutex.Unlock()

	client, ok := s.clientManager.clients[deviceID]
	if !ok {
		return
	}

	if room, ok := s.roomManager.rooms[client.RoomID]; ok {
		delete(room.Clients, client.ID)
	}

	close(client.Done)
	delete(s.clientManager.clients, deviceID)
}

// BroadcastToRoom sends a message to all clients in a room except the sender.
// It safely handles and removes unresponsive clients.
func (s *RoomService) BroadcastToRoom(roomID string, clipboardData *shared.ClipboardUpdate) error {
	// Lock both managers to ensure safe, concurrent access.
	s.roomManager.mutex.Lock()
	defer s.roomManager.mutex.Unlock()
	s.clientManager.mutex.Lock()
	defer s.clientManager.mutex.Unlock()

	room, ok := s.roomManager.rooms[roomID]
	if !ok {
		return fmt.Errorf("room %s not found", roomID)
	}

	var clientsToRemove []*types.Client
	for _, client := range room.Clients {
		if client.DeviceID == clipboardData.DeviceId {
			continue
		}
		select {
		case client.Send <- clipboardData:
		default:
			clientsToRemove = append(clientsToRemove, client)
		}
	}

	for _, client := range clientsToRemove {
		// Remove from room
		delete(room.Clients, client.ID)
		// Remove from system and signal shutdown
		if c, ok := s.clientManager.clients[client.DeviceID]; ok {
			close(c.Done)
			delete(s.clientManager.clients, client.DeviceID)
		}
	}

	return nil
}
