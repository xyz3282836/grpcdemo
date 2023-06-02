package main

import (
	"bytes"
	"crypto/tls"
	"encoding/binary"
	"fmt"
	v1 "grpcdemo/api/v1"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/gogo/protobuf/proto"
	"golang.org/x/net/http2"
	// "google.golang.org/protobuf/proto"
)

type GrpcClient struct {
	http.Client
}

func NewGrpcClient() GrpcClient {
	return GrpcClient{http.Client{Transport: plainTextH2Transport}}
}

var plainTextH2Transport = &http2.Transport{
	// 允许发起 http:// 请求
	AllowHTTP: true,
	// 默认 http2 会发起 tls 连接
	// 可以通过 DialTLS 拦截并改为发起普通 TCP 连接
	DialTLS: func(network, addr string, _ *tls.Config) (net.Conn, error) {
		return net.Dial(network, addr)
	},
}

func main() {
	c := http.Client{Transport: plainTextH2Transport}

	api := "http://127.0.0.1:9000/grpcdemo.v1.Hello/GetID"
	in := &v1.GetIdReq{}
	// in := &v1.HelloRequest{Name: "a"}
	// for i := 0; i < 1; i++ {
	// 	in.Name += "b"
	// }

	// 正式代码需要处理错误
	pb, _ := proto.Marshal(in)

	bs := make([]byte, 5)

	// 第一个字节默认为0，表示不压缩
	// 后四个字节以大端的形式保存 pb 消息长度
	binary.BigEndian.PutUint32(bs[1:], uint32(len(pb)))
	log.Printf("in  data len is %#b %#x %d lens is %d", bs, bs, binary.BigEndian.Uint32(bs[1:]), uint32(len(pb)))
	// 使用 MultiReader 「连接」两个 []byte 避免无谓的内存拷贝
	body := io.MultiReader(bytes.NewReader(bs), bytes.NewReader(pb))
	//body1 := io.MultiReader(bytes.NewReader(bs), bytes.NewReader(pb))
	//tmp, _ := io.ReadAll(body)
	//log.Printf("in  data is %#b", tmp)

	req, _ := http.NewRequest("POST", api, body)
	req.Header.Set("trailers", "TE")
	req.Header.Set("content-type", "application/grpc+proto")

	resp, _ := c.Do(req)
	defer func() {
		_ = resp.Body.Close()
	}()

	pb, _ = io.ReadAll(resp.Body)
	if status := resp.Header.Get("grpc-status"); status != "" {
		var c int
		var err error
		if c, err = strconv.Atoi(status); err != nil {
			return
		}
		err = fmt.Errorf("status %d err is %s", c, resp.Header.Get("grpc-message"))
		log.Fatalln(err)
		return
	}
	out := &v1.HelloReply{}
	err := proto.Unmarshal(pb[5:], out)
	//log.Printf("out data is %#b", pb)
	log.Printf("out data len is %#b %#x %d", pb[:5], pb[:5], binary.BigEndian.Uint32(pb[1:5]))

	// log.Printf("out pb is %#x err is %v", pb, err)
	log.Printf("err is %v", err)
	log.Printf("pb 6-7 %#b %d", pb[5:7], binary.BigEndian.Uint16(pb[5:7]))

	log.Printf("%s", string([]byte{0b1010}))

}
