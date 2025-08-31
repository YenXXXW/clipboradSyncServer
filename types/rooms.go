package types

import (
	"sync"

	"github.com/YenXXXW/clipboradSyncServer/shared"
)

type Room struct {
	ClipBoardContent string
	Clients          map[string]*Client
}

type Client struct {
	ID       string
	DeviceID string
	Conn     shared.StreamWriter
	RoomID   string
	Send     chan *shared.ClipboardUpdate
	Done     chan struct{}
	Mutex    sync.Mutex
}

type RoomService interface {
	CreateRoom() string
	JoinRoom(string, *Client) error
	RemoveFromRoom(string, string) error
	GetClient(string) (*Client, bool)
	GetRoom(string) (*Room, bool)
	CreateClient(string, string, shared.StreamWriter) *Client
	DeleteClient(string)
	BroadcastToRoom(string, *shared.ClipboardUpdate) error
}
