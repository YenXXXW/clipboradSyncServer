package service

import (
	"fmt"
	"sync"

	"github.com/YenXXXW/clipboradSyncServer/types"
	"github.com/google/uuid"
)

func NewRoom() *types.Room {
	return &types.Room{
		ClipBoardContent: "",
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

func (rm *RoomManager) LeaveRoom(roomID string, client *types.Client) error {
	rm.Mutex.Lock()
	defer rm.Mutex.Unlock()

	room, ok := rm.Rooms[roomID]
	if !ok {
		return fmt.Errorf("room with ID %s does not exist", roomID)
	}

	delete(room.Clients, client.ID)
	return nil
}
