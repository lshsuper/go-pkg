package jwt

type TokenRequest struct {
	Expire int
	SigningKey string
	Parms map[string]interface{}
}


type CheckTokenRequest struct {

	TokenStr string
	SigningKey string

}




