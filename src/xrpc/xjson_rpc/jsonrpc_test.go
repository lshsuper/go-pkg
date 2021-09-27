package xjson_rpc

import (
	"fmt"
	"testing"
	"time"
)

type UserService struct {
}

func (u *UserService) ServiceName() string {
	return "UserService"
}
func (u *UserService) GetUserName(req GetUserRequest, res *GetUserResponse) error {
	res.Name = "123456"
	return nil
}

type GetUserRequest struct {
}
type GetUserResponse struct {
	Name string
}

func TestDo(t *testing.T) {
	go func() {
		server, err := NewServer("127.0.0.1:10086")
		fmt.Println(err)
		uServer := new(UserService)
		server.Register(uServer)
		server.Start()
		fmt.Scanln()
	}()
	time.Sleep(time.Second * 3)
	go func() {
		for {

			client, _ := NewClient("127.0.0.1:10086")

			res := new(GetUserResponse)
			client.Do("UserService.GetUserName", GetUserRequest{}, res)
			client.Close()

			fmt.Println(res.Name)

		}
	}()

	time.Sleep(time.Second * 20)

}
