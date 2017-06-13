package auth

import "testing"

func TestAuthSuccess(t *testing.T) {

    tokenString, err := Auth("admin", "admin")
    if err != nil {
        t.Fatal(err)
    }

    userName, err := Validate(tokenString)
    if err != nil {
        t.Fatal(err)
    }

    if userName != "admin" {
        t.Fatal("validated userName is invalid!:", userName)
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
