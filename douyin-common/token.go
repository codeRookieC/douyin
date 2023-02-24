package douyin_common

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

func CreateToken(userName string) (tokenString string) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":      "lora-app-server",
		"aud":      "lora-app-server",
		"nbf":      time.Now().Unix(),
		"exp":      time.Now().Add(time.Hour).Unix(),
		"sub":      "userServer",
		"username": userName,
	})
	tokenString, err := token.SignedString([]byte("verysecret"))
	if err != nil {
		panic(err)
	}
	return tokenString
}
