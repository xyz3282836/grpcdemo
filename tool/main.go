package main

import (
	"context"
	"flag"
	"fmt"
	"strings"
	"time"

	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"

	"log"
	// https://github.com/gogo/protobuf/protoc-gen-gogo
	"github.com/golang/protobuf/proto"
	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"google.golang.org/grpc"
	rpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

type ErrorResponse struct {
	*rpb.ServerReflectionResponse_ErrorResponse
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("error_code:%d  error_message:%s", e.ErrorResponse.ErrorCode, e.ErrorResponse.ErrorMessage)
}

var kacp = keepalive.ClientParameters{
	Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
	Timeout:             time.Second,      // wait 1 second for ping ack before considering the connection dead
	PermitWithoutStream: false,            // send pings even without active streams
}

func main() {
	var ip string
	var port string
	flag.StringVar(&ip, "ip", "10.150.124.121", "-ip 127.0.0.1")
	flag.StringVar(&port, "port", "9000", "-port 9000")
	flag.Parse()

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithKeepaliveParams(kacp))
	// 建联grpc服务
	conn, err := grpc.Dial(ip+":"+port, opts...)
	if err != nil {
		log.Printf("grpc.Dial fail %v", err)
		panic("grpc dial fail")
	}
	defer conn.Close()

	reflectionList(context.TODO(), conn)
}

func reflectionList(ctx context.Context, conn *grpc.ClientConn) {
	// 获取grpc服务信息
	client := rpb.NewServerReflectionClient(conn)
	stream, err := client.ServerReflectionInfo(ctx, grpc.WaitForReady(true))
	if err != nil {
		log.Printf("ServerReflectionInfo err %v", err)
		return
	}
	services, err := listServices(ctx, stream)
	if err != nil {
		log.Printf("listServices err %v", err)
		return
	}
	var methods []string
	log.Printf("v2 api service len %d", len(services.ListServicesResponse.Service))
	for index, v := range services.ListServicesResponse.Service {
		log.Println()
		log.Printf("v2 service name is %s and index is %d", v.Name, index+1)
		info, e := getServiceInfo(ctx, stream, v.Name)
		if e != nil {
			log.Printf("service name %s 需要将它所在proto文件通过非gogo工具重新生成一次 err %v", v.Name, e)
			continue
		}
		for _, content := range info.FileDescriptorResponse.FileDescriptorProto {
			m := dpb.FileDescriptorProto{}
			if err = proto.Unmarshal(content, &m); err != nil {
				log.Printf("unmarshal file descriptor err %v service name is %s", err, m.Name)
				continue
			}
			for _, s := range m.Service {
				strArr := strings.Split(v.Name, ".")
				if strArr[len(strArr)-1] != s.GetName() {
					continue
				}
				log.Printf("-- v2 service name is %s and method len is %d", s.GetName(), len(s.Method))
				for _, method := range s.Method {
					methods = append(methods, fmt.Sprintf("/%s/%s", v.Name, method.GetName()))
					log.Printf("-- -- v2 api is /%s/%s", v.Name, method.GetName())
				}
			}
		}
	}
	log.Println()
	log.Printf("all api num is %d", len(methods))
	return
}

func listServices(ctx context.Context, stream rpb.ServerReflection_ServerReflectionInfoClient) (res *rpb.ServerReflectionResponse_ListServicesResponse, err error) {
	err = stream.Send(&rpb.ServerReflectionRequest{
		MessageRequest: &rpb.ServerReflectionRequest_ListServices{},
	})
	if err != nil {
		return
	}
	resp, err := stream.Recv()
	if err != nil {
		log.Printf("listServices err %v", err)
		return
	}
	res, ok := resp.MessageResponse.(*rpb.ServerReflectionResponse_ListServicesResponse)
	if !ok {
		log.Printf("listServices MessageResponse is not list")
		err = fmt.Errorf("listServices reponse not valid")
		return
	}
	return
}

func getServiceInfo(ctx context.Context, stream rpb.ServerReflection_ServerReflectionInfoClient, serviceName string) (res *rpb.ServerReflectionResponse_FileDescriptorResponse, err error) {
	err = stream.Send(&rpb.ServerReflectionRequest{
		MessageRequest: &rpb.ServerReflectionRequest_FileContainingSymbol{FileContainingSymbol: serviceName},
	})
	if err != nil {
		return
	}
	resp, err := stream.Recv()
	if err != nil {
		log.Printf("getServiceInfo err %v", err)
		return
	}
	res, ok := resp.MessageResponse.(*rpb.ServerReflectionResponse_FileDescriptorResponse)
	if !ok {
		err = fmt.Errorf("getServiceInfo reponse not valid")
		return
	}
	return
}
