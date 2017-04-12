package user

import (
    "testing"
    "time"
    "os"
)

var testUser = User{
    userId:     "0123456789",
    mac:        "aa:bb:cc:00:22:44",
    name:       "TestMan",
    lastAppear: time.Now(),
}

var ctx *DBContext

func setup() {
    println("setup")

    dbPath = "./users_test.db"
    os.Remove(dbPath)
    ctx = GetContext()
}

func teardown() {
    println("teardown")

    ctx.Close()
    os.Remove(dbPath)
}

func TestInsert(t *testing.T) {

    // insert
    err := ctx.InsertUser(testUser)
    if err != nil {
        t.Errorf("insert err: %v", err)
    }
}

func TestFind(t *testing.T) {

    // find
    found, err := ctx.FindUser(testUser.userId)
    if err != nil {
        t.Errorf("insert err: %v", err)
        return
    }

    if testUser.userId != found.userId ||
        testUser.mac != found.mac ||
        testUser.name != found.name ||
        !testUser.lastAppear.Equal(found.lastAppear) {
        t.Errorf("found user failed.\n expected: %v\n actual: %v\n",testUser, found)
    }
}

func TestUpdate(t *testing.T) {

    var userCopy = testUser

    // update
    userCopy.lastAppear = userCopy.lastAppear.Add(60)
    err := ctx.UpdateLastAppear(userCopy.userId, userCopy.lastAppear)
    if err != nil {
        t.Errorf("update err: %v", err)
    }
}

func TestErase(t *testing.T) {

    // erase
    err := ctx.EraseUser(testUser.userId)
    if err != nil {
        t.Errorf("erase err: %v", err)
    }
}

func TestMain(m *testing.M) {
    setup()
    ret := m.Run()
    if ret == 0 {
        teardown()
    }
    os.Exit(ret)
}