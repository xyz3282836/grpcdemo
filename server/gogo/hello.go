package server

import (
	"context"
	v1 "grpcdemo/api/v1/gogo"
	"math"

	"github.com/gogo/protobuf/types"
)

type HelloServer struct {
	// gRPC 要求每个服务的实现必须嵌入对应的结构体
	// 这个结构体也是自动生成的
	v1.HelloServer
}

func (s *HelloServer) SayHello(ctx context.Context, in *v1.HelloRequest) (*v1.HelloReply, error) {
	out := &v1.HelloReply{}
	out.Message = "Hello " + in.Name
	out.Num = int64(math.MaxInt32 + 1)

	return out, nil
}

func (s *HelloServer) GetID(ctx context.Context, in *v1.GetIdReq) (*v1.GetIDsResp, error) {
	out := &v1.GetIDsResp{}
	startId := 10000000
	for i := startId; i < (startId + 50000); i++ {
		out.Ids = append(out.Ids, int64(i))
	}

	return out, nil
}

func (s *HelloServer) GetView(ctx context.Context, in *v1.GetViewReq) (*v1.GetViewResp, error) {
	out := &v1.GetViewResp{}
	out.Name = in.GetName()
	req := in.GetViews()
	one := &v1.FansMedalOptions{}
	types.UnmarshalAny(req[0].GetOptions(), one)
	out.Num = one.GetUpMid()
	return out, nil
}
