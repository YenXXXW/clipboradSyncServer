package service

import (
	pb "github.com/YenXXXW/clipboradSyncServer/genproto/clipboardSync"
)

type StreamWriter interface {
	send(*pb.ClipboardContent) error
}
