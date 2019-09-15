package api

import (
	"log"

	"golang.org/x/net/context"
)

type Server struct {
}

func (s *Server) SayHello(ctx context.Context, in *PingMessage) (*PingMessage, error) {
	log.Printf("Received Message :%s", in.Greeting)
	return &PingMessage{Greeting: in.GetGreeting()}, nil
}
