package jwt

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
)


func GetToken(req TokenRequest)(string,error)  {
	signingKey := []byte(req.SigningKey)
	claims:=jwt.MapClaims{
		"exp":req.Expire,
	}

    for k,v:=range req.Parms{
		claims[k]=v
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	ss, err := token.SignedString(signingKey)
    return ss,err
}

//CheckToken  
func CheckToken(req CheckTokenRequest,out interface{})(err error)  {

   token,err:=jwt.Parse(req.TokenStr, func(t *jwt.Token) (interface{}, error) {
	   return []byte(req.SigningKey), nil
   })

   if err!=nil{
	   return err
   }

   err=token.Claims.Valid()
   if err!=nil{
	   return
   }

   //转化一下
   jByte,_:=json.Marshal(token.Claims)
   err=json.Unmarshal(jByte,out)
   return









}
