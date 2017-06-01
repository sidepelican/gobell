package udb

import (
    "time"
)

type User struct {
    UserId     string       `json:"user_id"`
    Mac        string       `json:"mac"`
    Name       string       `json:"name"`
    Note       string       `json:"note"`
    LastAppear time.Time    `json:"last_appear"`
}

type Users []User

func (users Users) Contains(userId string) bool {
    for _, u := range users {
        if u.UserId == userId {
            return true
        }
    }
    return false
}

func (users Users) Difference(minus Users) Users {
    ret := Users{}
    for _, u := range users {
        if minus.Contains(u.UserId) == false {
            ret = append(ret, u)
        }
    }
    return ret
}

func NewUser(userId string, mac string, name string) User {
    return User{
        UserId:     userId,
        Mac:        mac,
        Name:       name,
        Note:       "",
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
