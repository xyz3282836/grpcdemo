package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	v1 "grpcdemo/api/v1"
	"grpcdemo/server"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

var kasp = keepalive.ServerParameters{
	MaxConnectionIdle:     5000 * time.Second, // If a client is idle for 15 seconds, send a GOAWAY
	MaxConnectionAge:      6000 * time.Second, // If any connection is alive for more than 30 seconds, send a GOAWAY
	MaxConnectionAgeGrace: 5 * time.Second,    // Allow 5 seconds for pending RPCs to complete before forcibly closing connections
	Time:                  15 * time.Second,   // Ping the client if it is idle for 5 seconds to ensure the connection is still active
	Timeout:               1 * time.Second,    // Wait 1 second for the ping ack before assuming the connection is dead
}

func main() {
	ln, _ := net.Listen("tcp", "127.0.0.1:9000")

	// new grpc server
	grpcServer := grpc.NewServer(grpc.KeepaliveParams(kasp))
	// 注册接口
	v1.RegisterHelloServer(grpcServer, &server.HelloServer{})
	// 启动
	go grpcServer.Serve(ln)
	log.Printf("grpc server start")

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			grpcServer.GracefulStop()
			log.Printf("grpc server stop")
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}

}
