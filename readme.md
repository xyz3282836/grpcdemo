# 生成proto
`protoc --go_out=/Users/zhou/go/src/grpcdemo/ --go-grpc_out=/Users/zhou/go/src/grpcdemo api/v1/hello.proto`

# 启动grpc服务
`go run cmd/main.go`
## 查看grpc服务
`lsof -i :9000`
## 客户端请求
`go run client/main.go`
