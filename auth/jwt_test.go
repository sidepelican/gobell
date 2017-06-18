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
    orgTimeFunc := jwt.TimeFunc
    defer func() {
        jwt.TimeFunc = orgTimeFunc
    }()

    jwt.TimeFunc = func() (time.Time) {
        return time.Now().Add(- time.Hour * 24 * 30)
    }

    tokenString, err := Auth("admin", "admin")
    if err != nil {
        t.Fatal(err)
    }

    jwt.TimeFunc = time.Now

    info, err := Validate(tokenString)
    t.Log("tokenExpiresAt:", time.Unix(info.ExpiresAt, 0))
    t.Log("serverTime:", jwt.TimeFunc())

    if err == nil {
        t.Fatal("timeout error not handled")
    }
}

func TestUpdateExpires(t *testing.T) {
    orgTimeFunc := jwt.TimeFunc
    defer func() {
        jwt.TimeFunc = orgTimeFunc
    }()

    // yesterday created token.
    jwt.TimeFunc = func() (time.Time) {
        return time.Now().Add(- time.Hour * 24)
    }

    tokenString, err := Auth("admin", "admin")
    if err != nil {
        t.Fatal(err)
    }

    prevInfo, err := Validate(tokenString)
    if err != nil {
        t.Fatal("timeout error not handled")
    }

    // time is now
    jwt.TimeFunc = time.Now
    tokenString, err = UpdateExpires(tokenString)
    if err != nil {
        t.Fatal(err)
    }

    nowInfo, err := Validate(tokenString)
    if err != nil {
        t.Fatal(err)
    }

    // check
    if nowInfo.ExpiresAt < prevInfo.ExpiresAt {
        t.Log("prevExpiresAt:", time.Unix(prevInfo.ExpiresAt, 0))
        t.Log("nowExpiresAt:", time.Unix(nowInfo.ExpiresAt, 0))
        t.Fatal("failed to update Expiresat")
    }
}
