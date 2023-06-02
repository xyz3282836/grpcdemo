package main

import (
	"context"
	v1 "grpcdemo/grpcclient/protodiff/api/v1"
	"log"

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

	client := v1.NewHelloClient(conn)

	ret, err := client.SayHello(context.TODO(), &v1.HelloRequest{Name: "zhou"})
	if err != nil {
		log.Printf("client grpc fail %v", err)
		return
	}
	log.Printf("int64 num is %d", ret.Num)

}
