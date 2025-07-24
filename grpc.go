package main

import (
	"log"
)

type gRPCServer struct {
	addr string
}

func NewgRPCServer(addr string) *gRPCServer {
	return &gRPCServer{
		addr: addr,
	}
}

func (s *gRPCServer) Run() error {
}
