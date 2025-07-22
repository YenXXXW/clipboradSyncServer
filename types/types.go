package types

import (
	"context"
	"github.com/YenXXXW/clipboradSyncServer/genproto/clipboardSync"
	"google.golang.org/grpc"

	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type ClipboardSyncService interface {
	SubscribeClipBoardContentUpdate(*clipboardSync.SubscribeRequest, grpc.ServerStreamingServer[clipboardSync.ClipboardContent]) error
	SendClipBoardUpdate(context.Context, *clipboardSync.ClipboardUpdateRequest) (*emptypb.Empty, error)
}
