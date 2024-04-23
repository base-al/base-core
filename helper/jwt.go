package helper

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type CustomClaims struct {
	UserData interface{}
	jwt.StandardClaims
}

func Encode(userData interface{}, sk string) (string, error) {
	expireToken := time.Now().Add(time.Hour * 1).Unix()
	claims := CustomClaims{
		userData,
		jwt.StandardClaims{
			ExpiresAt: expireToken,
			Issuer:    "base.user",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(sk))
}

func Decode(tokenString string, sk string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(sk), nil
	})
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

func TokenValidation(token string, sk string) (resp *CustomClaims, err error) {
	claims, err := Decode(token, sk)
	if err != nil {
		return nil, errors.New("token not valid")
	}
	return claims, nil
}
