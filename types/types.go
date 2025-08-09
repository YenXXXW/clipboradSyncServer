package types

import (
	"context"

	//pb "github.com/YenXXXW/clipboradSyncServer/genproto/clipboardSync"
	"github.com/YenXXXW/clipboradSyncServer/shared"
)

type ClipboardSyncService interface {
	SubscribeClipBoardContentUpdate(string, string, shared.StreamWriter) error
	SendClipBoardUpdate(context.Context, string, *shared.ClipboardUpdate) error
}
