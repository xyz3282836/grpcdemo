package main

import (
	"google.golang.org/grpc"
	v1 "grpcdemo/api/v1"
	"grpcdemo/server"
	"net"
)

func main() {
	ln, _ := net.Listen("tcp", "127.0.0.1:9000")

	// new grpc server
	grpcServer := grpc.NewServer()
	// 注册接口
	v1.RegisterHelloServer(grpcServer, &server.HelloServer{})
	// 启动
	grpcServer.Serve(ln)

}
