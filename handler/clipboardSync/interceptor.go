package handler

import (
	"context"
	"fmt"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

var (
	errMissingMetadata = status.Errorf(codes.InvalidArgument, "missing metadata")
	errInvalidToken    = status.Errorf(codes.Unauthenticated, "invalid token")
)

func valid(authorization []string) bool {
	if len(authorization) < 1 {
		return false
	}
	token := strings.TrimPrefix(authorization[0], "Bearer ")
	// Perform the token validation here. For the sake of this example, the code
	// here forgoes any of the usual OAuth2 token validation and instead checks
	// for a token matching an arbitrary string.
	return token == "some-secret-token"
}

func logger(format string, a ...any) {
	fmt.Printf(format+"\n", a...)
}

func UnaryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	_, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errMissingMetadata
	}

	//authHeader := ""
	//if auths := md["authorization"]; len(auths) > 0 {
	//authHeader = auths[0]
	//}

	var clientIP string
	if p, ok := peer.FromContext(ctx); ok {
		// p.Addr is the network address of the client.
		// p.Addr.String() returns the address as a string.
		clientIP = p.Addr.String()
	} else {
		clientIP = "unknown"
	}

	logger("[%s] RPC called: %s | From: %s",
		time.Now().Format(time.RFC3339), // current date/time
		info.FullMethod,                 // gRPC service/method
		clientIP,                        // who called
	)

	//if !valid(md["authorization"]) {
	//return nil, errInvalidToken
	//}

	m, err := handler(ctx, req)
	if err != nil {
		logger("RPC failed with error: %v", err)
	}
	return m, err
}

type wrappedStream struct {
	grpc.ServerStream
}

func (w *wrappedStream) RecvMsg(m any) error {
	logger("Receive a message (Type: %T) at %s", m, time.Now().Format(time.RFC3339))
	return w.ServerStream.RecvMsg(m)
}

func (w *wrappedStream) SendMsg(m any) error {
	logger("Send a message (Type: %T) at %v", m, time.Now().Format(time.RFC3339))
	return w.ServerStream.SendMsg(m)
}

func newWrappedStream(s grpc.ServerStream) grpc.ServerStream {
	return &wrappedStream{s}
}

func StreamInterceptor(srv any, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	// authentication (token verification)
	_, ok := metadata.FromIncomingContext(ss.Context())
	if !ok {
		return errMissingMetadata
	}
	//if !valid(md["authorization"]) {
	//return errInvalidToken
	//}

	var clientIP string
	if p, ok := peer.FromContext(ss.Context()); ok {
		// p.Addr is the network address of the client.
		// p.Addr.String() returns the address as a string.
		clientIP = p.Addr.String()
	} else {
		clientIP = "unknown"
	}

	logger("[%s] RPC called: %s | From: %s",
		time.Now().Format(time.RFC3339), // current date/time
		info.FullMethod,                 // gRPC service/method
		clientIP,                        // who called
	)

	err := handler(srv, newWrappedStream(ss))
	if err != nil {
		logger("RPC failed with error: %v", err)
	}
	return err
}
