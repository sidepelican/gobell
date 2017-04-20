package udb

import (
    "time"
)

type User struct {
    UserId     string
    Mac        string
    Name       string
    LastAppear time.Time
}

type Users []User

func NewUser(userId string, mac string, name string) User {
    return User{
        UserId:     userId,
        Mac:        mac,
        Name:       name,
        LastAppear: time.Now(),
    }
}

func (u Users) Len() int {
    return len(u)
}

func (u Users) Swap(i, j int) {
    u[i], u[j] = u[j], u[i]
}

func (u Users) Less(i, j int) bool {
    return u[i].Name < u[j].Name
}