package handler

import (
	pb "github.com/YenXXXW/clipboradSyncServer/genproto/clipboardSync"
	"github.com/YenXXXW/clipboradSyncServer/shared"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type grpcStreamWrapper struct {
	stream grpc.ServerStreamingServer[pb.UpdateEvent]
}

func (g *grpcStreamWrapper) Send(update *shared.UpdateEvent) error {
	var grpcUpdate *pb.UpdateEvent

	// Check which type is non-nil
	if update.ClipboardUpdate != nil {
		grpcUpdate = &pb.UpdateEvent{
			Event: &pb.UpdateEvent_ClipboardUpdate{
				ClipboardUpdate: &pb.ClipboardUpdate{
					DeviceId: update.ClipboardUpdate.DeviceId,
					Content: &pb.ClipboardContent{
						Text: update.ClipboardUpdate.Content.Text,
					},
				},
			},
		}
	} else if update.ValidateJoin != nil {
		grpcUpdate = &pb.UpdateEvent{
			Event: &pb.UpdateEvent_ValidateJoin{
				ValidateJoin: &pb.ValidateJoin{
					ValidateRoom: &pb.Validate{
						Success: update.ValidateJoin.ValidateRoom.Success,
						Message: update.ValidateJoin.ValidateRoom.Message,
					},
					CheckClient: &pb.Validate{
						Success: update.ValidateJoin.CheckClient.Success,
						Message: update.ValidateJoin.CheckClient.Message,
					},
				},
			},
		}
	} else {
		return nil
	}

	return g.stream.Send(grpcUpdate)
}

func (g *grpcStreamWrapper) Context() context.Context {
	return g.stream.Context()
}
