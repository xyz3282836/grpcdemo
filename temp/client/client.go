package main

import (
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

var addr = flag.String("addr", "localhost:50052", "the address to connect to")

var kacp = keepalive.ClientParameters{
	Time:                5 * time.Second, // send pings every 10 seconds if there is no activity
	Timeout:             time.Second,     // wait 1 second for ping ack before considering the connection dead
	PermitWithoutStream: false,           // send pings even without active streams
}

func main() {
	flag.Parse()

	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithKeepaliveParams(kacp))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// c := pb.NewEchoClient(conn)

	// ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	// defer cancel()
	// fmt.Println("Performing unary request")
	// res, err := c.UnaryEcho(ctx, &pb.EchoRequest{Message: "keepalive demo"})
	// if err != nil {
	// 	log.Fatalf("unexpected error from UnaryEcho: %v", err)
	// }
	// fmt.Println("RPC response:", res)
	select {} // Block forever; run with GODEBUG=http2debug=2 to observe ping frames and GOAWAYs due to idleness.
}
