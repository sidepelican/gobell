package auth

import (
    "testing"
    "time"

    jwt "github.com/dgrijalva/jwt-go"
)

func TestAuthSuccess(t *testing.T) {

    tokenString, err := Auth("admin", "admin")
    if err != nil {
        t.Fatal(err)
    }

    info, err := Validate(tokenString)
    if err != nil {
        t.Fatal(err)
    }

    if info.Name != "admin" {
        t.Fatal("validated userName is invalid!:", info.Name)
    }
}

func TestAuthFail(t *testing.T) {

    tokenString, err := Auth("admin", "1234")
    if err == nil {
        t.Fatal()
    }
    if tokenString != "" {
        t.Fatal("tokenString is not blank:", tokenString)
    }
}

func TestTimeout(t *testing.T) {

    tokenString, err := Auth("admin", "admin")
    if err != nil {
        t.Fatal(err)
    }

    orgTimeFunc := jwt.TimeFunc
    jwt.TimeFunc = func() (time.Time) {
        return time.Now().Add(time.Hour * 24 * 30)
    }

    info, err := Validate(tokenString)
    t.Log(time.Unix(info.ExpiresAt, 0))
    t.Log(jwt.TimeFunc())

    if err == nil {
        t.Fatal("timeout error not handled")
    }

    jwt.TimeFunc = orgTimeFunc
}
