package xjson_rpc

import (
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

//IJsonRpcService jsonRpc-service-interface
type IJsonRpcService interface {
	ServiceName() string
}

type jsonrpcServer struct {
	server net.Listener
}

//NewServer 实例化一个服务
func NewServer(addr string) (*jsonrpcServer, error) {

	server := new(jsonrpcServer)
	lis, err := net.Listen("tcp", addr)

	if err != nil {
		return nil, err
	}

	server.server = lis
	return server, nil

}

//Register 注册服务
func (server *jsonrpcServer) Register(service IJsonRpcService) error {
	err := rpc.RegisterName(service.ServiceName(), service)
	return err
}

//Start 开启服务
func (server *jsonrpcServer) Start() {

	for {

		conn, err := server.server.Accept()
		if err != nil {
			continue
		}

		go func(conn net.Conn) { // 并发处理客户端请求
			jsonrpc.ServeConn(conn)
		}(conn)

	}

}
