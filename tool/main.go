package main

import (
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	grpcrf "github.com/jhump/protoreflect/grpcreflect"
)

var kacp = keepalive.ClientParameters{
	Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
	Timeout:             time.Second,      // wait 1 second for ping ack before considering the connection dead
	PermitWithoutStream: false,            // send pings even without active streams
}

func main() {
	var addr string
	flag.StringVar(&addr, "addr", "10.150.134.199:9000", "-addr 127.0.0.1:9000")
	flag.Parse()

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithKeepaliveParams(kacp))
	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		log.Printf("grpc.Dial fail %v", err)
		panic("grpc dial fail")
	}
	defer conn.Close()

	TestGetRelectApiList(conn)
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
		log.Printf("service name is %s", e)
		d, err := c.ResolveService(e)
		if err != nil {
			log.Printf("ResolveService err %v", err)
			return
		}
		for _, ee := range d.GetMethods() {
			log.Printf("api is /%s/%s", e, ee.GetName())
		}
	}
}
