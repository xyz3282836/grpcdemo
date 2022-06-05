# 生成proto
有任何改动执行，绝对路径，换成自己机器的
`protoc --gofast_out=plugins=grpc:/Users/zhou/go/src/grpcdemo/ api/v1/hello.proto`
`protoc --gogofast_out=plugins=grpc:/Users/zhou/go/src/grpcdemo/ api/v1/hello.proto`

# 启动grpc服务
启动grpc服务
`go run cmd/main.go`

## 查看grpc服务
查看grpc服务进程
`lsof -i :9000`

## 客户端请求
模拟grpc请求，以http的方式
`go run client/main.go`
