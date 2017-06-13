package auth

import (
    "time"
    "errors"

    jwt "github.com/dgrijalva/jwt-go"
)

var secret = []byte("test_secret")

type authInfo struct {
    jwt.StandardClaims
    Name string     `json:"name"`
    Time time.Time  `json:"time"`
}

func Auth(name string, pass string) (string, error) {

    if name == "admin" && pass == "admin" {
        token := jwt.NewWithClaims(jwt.SigningMethodHS256, authInfo{
            Name: name,
            Time: time.Now(),
        })

        return token.SignedString(secret)
    }

    return "", errors.New("failed to auth")
}

func Validate(tokenString string) (string, error) {

    info := authInfo{}
    token, err := jwt.ParseWithClaims(tokenString, &info, func(token *jwt.Token) (interface{}, error){
        return secret, nil
    })
    if err != nil {
        return "", err
    }
    if !token.Valid {
        return "", errors.New("token.Valid is false. something wrong.")
    }

    return info.Name, nil
}