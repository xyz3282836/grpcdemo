package main

import (
	"context"
	"encoding/json"
	v1 "grpcdemo/api/v1"
	"log"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial("127.0.0.1:9000", opts...)
	if err != nil {
		log.Printf("grpc.Dial fail %v", err)
		panic("grpc dial fail")
	}
	defer conn.Close()

	TestSayHello(conn)
	TestGetView(conn)
}

func TestSayHello(conn *grpc.ClientConn) {
	client := v1.NewHelloClient(conn)

	ret, err := client.SayHello(context.TODO(), &v1.HelloRequest{Name: "zhou"})
	if err != nil {
		log.Printf("client grpc fail %v", err)
		return
	}
	log.Printf("ret is %v", ret)
}

func TestGetView(conn *grpc.ClientConn) {
	client := v1.NewHelloClient(conn)

	data := &v1.FansMedalOptions{UpMid: int64(112233)}
	anydata, _ := ptypes.MarshalAny(data)

	req := &v1.GetViewReq{Name: "zhou", Views: []*v1.UserViewItem{{View: v1.UserViewEnum_BASE_INFO_VIEW, Options: anydata}}}

	reqstr, _ := json.Marshal(req)

	ret, err := client.GetView(context.TODO(), req)
	if err != nil {
		log.Printf("client grpc fail %v", err)
		return
	}
	log.Printf("req is %s\n", reqstr)
	log.Printf("ret is %v", ret)
}
