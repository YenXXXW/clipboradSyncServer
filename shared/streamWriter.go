package shared

import (
	"context"

	pb "github.com/YenXXXW/clipboradSyncServer/genproto/clipboardSync"
)

type StreamWriter interface {
	Send(*pb.ClipboardContent) error
	Context() context.Context
}
