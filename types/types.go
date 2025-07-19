package types

import (
	"context"
	"github.com/YenXXXW/clipboradSyncServer/genproto/clipboardSync"
	"google.golang.org/grpc"

	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

type ClipboardSyncService interface {
	SubscribeClipBoardContentUpdate(*clipboardSync.SubscribeRequest, grpc.ServerStreamingServer[clipboardSync.ClipBoardContent]) error
	SendClipBoardUpdate(context.Context, *clipboardSync.ClipBoardContent) (*emptypb.Empty, error)
}
