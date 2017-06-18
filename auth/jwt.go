package auth

import (
    "time"
    "errors"

    jwt "github.com/dgrijalva/jwt-go"
)

var secret = []byte("test_secret")

func generateNewExpiredTime() time.Time {
    const loginSessionTimeOut = time.Hour * 24 * 7
    return jwt.TimeFunc().Add(loginSessionTimeOut)
}

type AuthInfo struct {
    jwt.StandardClaims
    Name string     `json:"name"`
}

func Auth(name string, pass string) (string, error) {

    if name == "admin" && pass == "admin" {
        token := jwt.NewWithClaims(jwt.SigningMethodHS256, AuthInfo{
            jwt.StandardClaims{
                ExpiresAt: generateNewExpiredTime().Unix(),
            },
            name,
        })

        return token.SignedString(secret)
    }

    return "", errors.New("failed to auth")
}

func Validate(tokenString string) (info AuthInfo, err error) {

    _, err = jwt.ParseWithClaims(tokenString, &info, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("Unexpected siging method")
        }
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

func UpdateExpires(tokenString string) (string, error) {

    info, err := Validate(tokenString)
    if err != nil {
        return "", err
    }

    info.ExpiresAt = generateNewExpiredTime().Unix()

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, info)
    return token.SignedString(secret)
}
