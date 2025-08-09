package shared

import (
	"context"
)

type ClipboardContent struct {
	Text string
}

type ClipboardUpdate struct {
	DeviceId string
	Content  ClipboardContent
}

type StreamWriter interface {
	Send(*ClipboardUpdate) error
	Context() context.Context
}
