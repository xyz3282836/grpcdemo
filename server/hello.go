package server

import (
	"context"
	"fmt"
	v1 "grpcdemo/api/v1"
	"log"
	"math"
	"time"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
	gm "google.golang.org/grpc/metadata"
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
	out.Num = int64(0)
	grpc.SetTrailer(ctx, gm.Pairs([]string{"ccppuu", "lalal"}...))
	grpc.SetHeader(ctx, gm.Pairs([]string{"color", "red"}...))
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
	_ = ptypes.UnmarshalAny(req[0].GetOptions(), one)
	out.Num = one.GetUpMid()
	return out, nil
}

func (s *HelloServer) GetStream(req *v1.GetStreamReq, srv v1.Hello_GetStreamServer) (err error) {
	reqName := req.GetName()
	for i := 1; i < 11; i++ {
		//if i == 4 {
		//	return fmt.Errorf("build err")
		//}
		err = srv.Send(&v1.GetStreamResp{
			Result: fmt.Sprintf("stream return no.%d %d and resp time %s", i, time.Now().Unix(), reqName),
		})
		time.Sleep(2 * time.Second)
		if err != nil {
			log.Printf("stream err %v", err)
			return
		}
	}
	return
}

var xxx = "dsfsdfdsfsdfsdfsdfsdfdsfsddfsdsfsdfdsfsdfsdfsdfsdfdsfsddfsdsfsdfdsfsdfsdfsdfsdfdsfsddfsdsfsdfdsfsdfsdfsdfsdfdsfsddfsdsfsdfdsfsdfsdfsdfsdfdsfsddfsdsfsdfdsfsdfsdfsdfsdfdsfsddfsdsfsdfdsfsdfsdfsdfsdfdsfsddfsdsfsdfdsfsdfsdfsdfsdfdsfsddfsdsfsdfdsfsdfsdfsdfsdfdsfsddfsdsfsdfdsfsdfsdfsdfsdfdsfsddfsdsfsdfdsfsdfsdfsdfsdfdsfsddfsdsfsdfdsfsdfsdfsdfsdfdsfsddfsdsfsdfdsfsdfsdfsdfsdfdsfsddfsdsfsdfdsfsdfsdfsdfsdfdsfsddfsdsfsdfdsfsdfsdfsdfsdfdsfsddfsdsfsdfdsfsdfsdfsdfsdfdsfsddfsdsfsdfdsfsdfsdfsdfsdfdsfsddfsdsfsdfdsfsdfsdfsdfsdfdsfsddfsdsfsdfdsfsdfsdfsdfsdfdsfsddfsdsfsdfdsfsdfsdfsdfsdfdsfsddfsdsfsdfdsfsdfsdfsdfsdfdsfsddfsdsfsdfdsfsdfsdfsdfsdfdsfsddfsdsfsdfdsfsdfsdfsdfsdfdsfsddfsdsfsdfdsfsdfsdfsdfsdfdsfsddfsdsfsdfdsfsdfsdfsdfsdfdsfsddfsdsfsdfdsfsdfsdfsdfsdfdsfsddfsdsfsdfdsfsdfsdfsdfsdfdsfsddfsdsfsdfdsfsdfsdfsdfsdfdsfsddfsdsfsdfdsfsdfsdfsdfsdfdsfsddfsdsfsdfdsfsdfsdfsdfsdfdsfsddfsdsfsdfdsfsdfsdfsdfsdfdsfsddfsdsfsdfdsfsdfsdfsdfsdfdsfsddfsdsfsdfdsfsdfsdfsdfsdfdsfsddfsdsfsdfdsfsdfsdfsdfsdfdsfsddfsdsfsdfdsfsdfsdfsdfsdfdsfsddfsdsfsdfdsfsdfsdfsdfsdfdsfsddfsdsfsdfdsfsdfsdfsdfsdfdsfsddfsdsfsdfdsfsdfsdfsdfsdfdsfsddfsdsfsdfdsfsdfsdfsdfsdfdsfsddfsdsfsdfdsfsdfsdfsdfsdfdsfsddfsdsfsdfdsfsdfsdfsdfsdfdsfsddfsdsfsdfdsfsdfsdfsdfsdfdsfsddfsdsfsdfdsfsdfsdfsdfsdfdsfsddfsdsfsdfdsfsdfsdfsdfsdfdsfsddfsdsfsdfdsfsdfsdfsdfsdfdsfsddfsdsfsdfdsfsdfsdfsdfsdfdsfsddfsdsfsdfdsfsdfsdfsdfsdfdsfsddfsdsfsdfdsfsdfsdfsdfsdfdsfsddfsdsfsdfdsfsdfsdfsdfsdfdsfsddfsdsfsdfdsfsdfsdfsdfsdfdsfsddfsdsfsdfdsfsdfsdfsdfsdfdsfsddfsdsfsdfdsfsdfsdfsdfsdfdsfsddfsdsfsdfdsfsdfsdfsdfsdfdsfsddfsdsfsdfdsfsdfsdfsdfsdfdsfsddfsdsfsdfdsfsdfsdfsdfsdfdsfsddfsdsfsdfdsfsdfsdfsdfsdfdsfsddfsdsfsdfdsfsdfsdfsdfsdfdsfsddfsdsfsdfdsfsdfsdfsdfsdfdsfsddfsdsfsdfdsfsdfsdfsdfsdfdsfsddfsdsfsdfdsfsdfsdfsdfsdfdsfsddfsdsfsdfdsfsdfsdfsdfsdfdsfsddfsdsfsdfdsfsdfsdfsdfsdfdsfsddfsdsfsdfdsfsdfsdfsdfsdfdsfsddfs"
