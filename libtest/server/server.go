package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/features/proto/echo"
	"google.golang.org/grpc/keepalive"
)

var port = flag.Int("port", 50052, "port number")

var kaep = keepalive.EnforcementPolicy{
	MinTime:             10 * time.Second, // If a client pings more than once every 5 seconds, terminate the connection
	PermitWithoutStream: false,             // Allow pings even when there are no active streams
}

var kasp = keepalive.ServerParameters{
	MaxConnectionIdle:     5000 * time.Second, // If a client is idle for 15 seconds, send a GOAWAY
	MaxConnectionAge:      6000 * time.Second, // If any connection is alive for more than 30 seconds, send a GOAWAY
	MaxConnectionAgeGrace: 5 * time.Second,    // Allow 5 seconds for pending RPCs to complete before forcibly closing connections
	Time:                  15 * time.Second,    // Ping the client if it is idle for 5 seconds to ensure the connection is still active
	Timeout:               1 * time.Second,    // Wait 1 second for the ping ack before assuming the connection is dead
}

// server implements EchoServer.
type server struct {
	pb.UnimplementedEchoServer
}

func (s *server) UnaryEcho(ctx context.Context, req *pb.EchoRequest) (*pb.EchoResponse, error) {
	return &pb.EchoResponse{Message: req.Message}, nil
}

func main() {
	flag.Parse()

	address := fmt.Sprintf(":%v", *port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(grpc.KeepaliveEnforcementPolicy(kaep), grpc.KeepaliveParams(kasp))
	pb.RegisterEchoServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
