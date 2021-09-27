package xjson_rpc

import (
	"net/rpc"
	"net/rpc/jsonrpc"
)

type jsonRpcClient struct {
	client *rpc.Client
}

//NewClient 实例化一个客户端
func NewClient(addr string) (*jsonRpcClient, error) {

	client, err := jsonrpc.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &jsonRpcClient{
		client: client,
	}, nil

}

//Close 关闭
func (c *jsonRpcClient) Close() {

	c.client.Close()
	return
}

//Do 执行
func (c *jsonRpcClient) Do(method string, req interface{}, res interface{}) error {

	err := c.client.Call(method, req, res)
	return err
}
