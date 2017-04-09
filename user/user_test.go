package user

import (
    "testing"
    "time"
)

func TestInsert(t *testing.T) {

    dbPath = "./users_test.db"

    // insert
    user := User{
        userId:     "0123456789",
        mac:        "aa:bb:cc:00:22:44",
        lastAppear: time.Now(),
    }
    err := InsertUser(user)
    if err != nil {
        t.Errorf("insert err: %v", err)
    }
}

func TestUpdate(t *testing.T) {

    dbPath = "./users_test.db"

    user := User{
        userId:     "0123456789",
        mac:        "aa:bb:cc:00:22:44",
        lastAppear: time.Now(),
    }

    // update
    user.lastAppear = user.lastAppear.Add(60)
    err := UpdateLastAppear(user.userId, user.lastAppear)
    if err != nil {
        t.Errorf("update err: %v", err)
    }
}

func TestFind(t *testing.T) {

    dbPath = "./users_test.db"

    user := User{
        userId:     "0123456789",
        mac:        "aa:bb:cc:00:22:44",
        lastAppear: time.Now(),
    }

    // find
    found, err := FindUser("0123456789")
    if err != nil {
        t.Errorf("insert err: %v", err)
    }
    if user.userId != found.userId ||
        user.mac != found.mac ||
        user.lastAppear != found.lastAppear {
        t.Error("found user failed")
    }
}
