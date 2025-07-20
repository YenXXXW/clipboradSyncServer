package types

import "google.golang.org/grpc"

type Room struct {
	ClipBoardContent string
	Clients          map[string]*Client
}

type Client struct {
	ID     string
	Conn   grpc.ServerStream
	RoomID string
	Send   chan []byte
}

type RoomService interface {
	CreateRoom() string
	JoinRoom(string, *Client) error
	LeaveRoom(string, *Client) error
}
