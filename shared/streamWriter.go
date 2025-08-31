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

type Validate struct {
	Success bool
	Message string
}

type ValidateJoin struct {
	ValidateRoom Validate
	CheckClient  Validate
}

type UpdateEvent struct {
	ClipboardUpdate *ClipboardUpdate
	ValidateJoin    *ValidateJoin
}

type StreamWriter interface {
	Send(*UpdateEvent) error
	Context() context.Context
}
