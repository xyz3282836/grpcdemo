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
	// server 发送给client的header和trailer
	grpc.SetTrailer(ctx, gm.Pairs("ccppuu", "99"))
	grpc.SetTrailer(ctx, gm.Pairs("mem", "4"))
	grpc.SetHeader(ctx, gm.Pairs("color", "red"))
	// 只能调用一次，在这个之下再去set header就无效了
	grpc.SendHeader(ctx, gm.Pairs("final", "send"))
	grpc.SetHeader(ctx, gm.Pairs("cluster", "sh"))
	// server 读取client的metadata
	md, ok := gm.FromIncomingContext(ctx)
	log.Printf("md is %v ok %v", md, ok)

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
