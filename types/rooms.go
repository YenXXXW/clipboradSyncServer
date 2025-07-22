package types

import (
	"google.golang.org/grpc"

	pb "github.com/YenXXXW/clipboradSyncServer/genproto/clipboardSync"
)

type Room struct {
	ClipBoardContent string
	Clients          map[string]*Client
}

type Client struct {
	ID     string
	Conn   grpc.ServerStream
	RoomID string
	Send   chan *pb.ClipboardContent
}

type RoomService interface {
	CreateRoom() string
	JoinRoom(string, *Client) error
	RemoveFromRoom(string, *Client) error
}
