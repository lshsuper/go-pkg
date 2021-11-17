package main

import (
	"fmt"
	"github.com/lshsuper/go-pkg/src/oauth/jwt"
)


type UserInfo struct {
	UserID int `json:"userId"`
	Name string  `json:"name"`
}

func main() {




	token,err:=jwt.GetToken(jwt.TokenRequest{
		SigningKey: "abcdefg",
		Parms: map[string]interface{}{
			"userId":123456,
			"name":"lsh",
		},
	})


	fmt.Println(token,err)

	  m:=&UserInfo{}

	err=jwt.CheckToken(jwt.CheckTokenRequest{
		TokenStr: token,
		SigningKey:"abcdefg",
	},m)




	fmt.Println(m,err)



}
