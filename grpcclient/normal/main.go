package main

import (
	"context"
	"encoding/json"
	"fmt"
	v1 "grpcdemo/api/v1"
	"io"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes"
	grpcrf "github.com/jhump/protoreflect/grpcreflect"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	gm "google.golang.org/grpc/metadata"
)

var kacp = keepalive.ClientParameters{
	Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
	Timeout:             time.Second,      // wait 1 second for ping ack before considering the connection dead
	PermitWithoutStream: false,            // send pings even without active streams
}

type myClientInterceptor struct{}

func (i *myClientInterceptor) UnaryClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	// 调用 gRPC 方法之前的处理逻辑
	err := invoker(ctx, method, req, reply, cc, opts...)

	// 获取 trailer 元数据
	trailer, ok := gm.FromOutgoingContext(ctx)
	fmt.Println("Trailer:ok", ok)
	fmt.Println("Trailer:", trailer.Get("ccppuu"))

	return err
}

func main() {
	var opts []grpc.DialOption
	// mci:=&myClientInterceptor{}
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithKeepaliveParams(kacp))

	conn, err := grpc.Dial("127.0.0.1:9000", opts...)
	if err != nil {
		log.Printf("grpc.Dial fail %v", err)
		panic("grpc dial fail")
	}
	defer conn.Close()

	TestSayHello(conn)
	//TestGetView(conn)
	// TestStream(conn)
	// TestGetRelectApiList(conn)
	// select {}
}

func TestSayHello(conn *grpc.ClientConn) {
	client := v1.NewHelloClient(conn)
	ctx := context.Background()
	var header, trailer gm.MD
	ret, err := client.SayHello(ctx, &v1.HelloRequest{Name: "zhou"}, grpc.Header(&header), grpc.Trailer(&trailer))
	if err != nil {
		log.Printf("client grpc fail %v", err)
		return
	}

	log.Printf("ret is %v header is %v trailer is %v", ret, header, trailer)
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

func TestGetRelectApiList(conn *grpc.ClientConn) {
	c := grpcrf.NewClientAuto(context.TODO(), conn)
	slist, err := c.ListServices()
	if err != nil {
		log.Printf("err %v", err)
		return
	}
	log.Printf("api len %d", len(slist))
	for _, e := range slist {
		d, _ := c.ResolveService(e)
		for _, ee := range d.GetMethods() {
			log.Printf("api is /%s/%s", e, ee.GetName())
		}
	}
}

func TestStream(conn *grpc.ClientConn) {
	client := v1.NewHelloClient(conn)

	req := &v1.GetStreamReq{Name: fmt.Sprintf("req time %d", time.Now().Unix())}

	cli, err := client.GetStream(context.TODO(), req)
	if err != nil {
		log.Printf("first init cli err %v", err)
		return
	}
	retry := true
	var tryTime int
	var recvErr error
	for tryTime < 5 && retry {
		tryTime++
		for {
			ret := &v1.GetStreamResp{}
			ret, recvErr = cli.Recv()
			if recvErr != nil {
				if recvErr == io.EOF {
					recvErr = nil
					break
				}
				log.Printf("recv err %v", recvErr)
				break
			}
			log.Printf("recv is %s", ret.GetResult())
		}
		if recvErr != nil {
			retry = true
			cli, err = client.GetStream(context.TODO(), req)
			if err != nil {
				log.Printf("cli err %v", err)
				return
			}
			log.Printf("try time %d", tryTime)
			time.Sleep(time.Second)
		} else {
			retry = false
		}
	}
	if retry != false {
		log.Printf("final err %v", recvErr)
		return
	}
	return
}
