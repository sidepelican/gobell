package auth

import (
    "time"
    "errors"

    jwt "github.com/dgrijalva/jwt-go"
)

var secret = []byte("test_secret")
const loginSessionTimeOut = time.Hour * 24 * 7

type AuthInfo struct {
    jwt.StandardClaims
    Name string     `json:"name"`
}

func Auth(name string, pass string) (string, error) {

    if name == "admin" && pass == "admin" {
        token := jwt.NewWithClaims(jwt.SigningMethodHS256, AuthInfo{
            jwt.StandardClaims {
                ExpiresAt: time.Now().Add(loginSessionTimeOut).Unix(),
            },
            name,
        })

        return token.SignedString(secret)
    }

    return "", errors.New("failed to auth")
}

func Validate(tokenString string) (info AuthInfo, err error) {

    _, err = jwt.ParseWithClaims(tokenString, &info, func(token *jwt.Token) (interface{}, error){
        return secret, nil
    })
    if err != nil {
        return
    }
    if err = info.Valid(); err != nil {
        return
    }

    return
}