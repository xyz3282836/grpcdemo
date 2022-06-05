package server

import (
	context "context"
	"grpcdemo/api/v1"
)

type HelloServer struct {
	// gRPC 要求每个服务的实现必须嵌入对应的结构体
	// 这个结构体也是自动生成的
	v1.UnimplementedHelloServer
}

func (s *HelloServer) SayHello(ctx context.Context, in *v1.HelloRequest) (*v1.HelloReply, error) {
	out := &v1.HelloReply{}
	out.Message = "Hello " + in.Name

	return out, nil
}
