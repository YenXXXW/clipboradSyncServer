package types

import (
	pb "github.com/YenXXXW/clipboradSyncServer/genproto/clipboardSync"
	"github.com/YenXXXW/clipboradSyncServer/shared"
	//"github.com/YenXXXW/clipboradSyncServer/shared"
)

type Room struct {
	ClipBoardContent string
	Clients          map[string]*Client
}

type Client struct {
	ID     string
	Conn   shared.StreamWriter
	RoomID string
	Send   chan *pb.ClipboardContent
}

type RoomService interface {
	CreateRoom() string
	JoinRoom(string, *Client) error
	RemoveFromRoom(string, string) error
	GetClient(string) (*Client, bool)
	CreateClient(string, string, shared.StreamWriter) *Client
	DeleteClient(string)
	BroadcastToRoom(string, *pb.ClipboardContent) error
}
