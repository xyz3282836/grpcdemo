# 生成proto

有任何改动执行，绝对路径，换成自己机器的

## go

```sh
protoc --go_out=/Users/zhou/go/src/grpcdemo/ --go-grpc_out=/Users/zhou/go/src/grpcdemo/ api/v1/hello.proto
```

### client diffproto

```sh
protoc --go_out=/Users/zhou/go/src/grpcdemo/grpcclient/protodiff --go-grpc_out=/Users/zhou/go/src/grpcdemo/grpcclient/protodiff grpcclient/protodiff/api/v1/hello.proto
```

## gogo

```sh
protoc --gofast_out=plugins=grpc,Mgoogle/protobuf/any.proto=github.com/gogo/protobuf/types:/Users/zhou/go/src/grpcdemo/ api/v1/gogo/hello.proto
protoc --gogofast_out=plugins=grpc,Mgoogle/protobuf/any.proto=github.com/gogo/protobuf/types:/Users/zhou/go/src/grpcdemo/ api/v1/gogo/hello.proto
```

## 启动grpc服务

启动grpc服务
`go run cmd/gogo/main.go`
`go run cmd/main.go`

## 查看grpc服务

查看grpc服务进程
`lsof -i :9000`

## 客户端请求

模拟grpc请求，以http的方式
`go run httpclient/main.go`
grpc客户端请求
`go run grpcclient/normal/gogo/main.go`
`go run grpcclient/normal/main.go`

## 扫描grpc server api list

`go run main.go -ip 10.150.255.153`

## grpc unary header trailer

```sh
go run cmd/main.go
go run grpcclient/normal/main.go
```
