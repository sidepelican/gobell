package udb

import (
    "testing"
    "time"
    "os"
    "strings"
)

var testUser = User{
    UserId:     "0123456789",
    Mac:        "aa:bb:cc:00:22:44",
    Name:       "TestMan",
    Note:       "test note 🥓 test 😀",
    LastAppear: time.Now(),
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
    found, err := ctx.FindUser(testUser.UserId)
    if err != nil {
        t.Errorf("insert err: %v", err)
        return
    }

    if testUser.UserId != found.UserId ||
        testUser.Mac != found.Mac ||
        testUser.Name != found.Name ||
        !testUser.LastAppear.Equal(found.LastAppear) {
        t.Errorf("found user failed.\n expected: %v\n actual: %v\n", testUser, found)
    }

    // mac
    found, err = ctx.FindMac(testUser.Mac)
    if err != nil {
        t.Errorf("insert err: %v", err)
        return
    }

    if testUser.UserId != found.UserId ||
        testUser.Mac != found.Mac ||
        testUser.Name != found.Name ||
        !testUser.LastAppear.Equal(found.LastAppear) {
        t.Errorf("found user failed.\n expected: %v\n actual: %v\n", testUser, found)
    }

    // mac uppercase
    found, err = ctx.FindMac(strings.ToUpper(testUser.Mac))
    if err != nil {
        t.Errorf("insert err: %v", err)
        return
    }

    if testUser.UserId != found.UserId ||
        testUser.Mac != found.Mac ||
        testUser.Name != found.Name ||
        !testUser.LastAppear.Equal(found.LastAppear) {
        t.Errorf("found user failed.\n expected: %v\n actual: %v\n", testUser, found)
    }
}

func TestUpdate(t *testing.T) {

    var userCopy = testUser

    // update
    userCopy.LastAppear = testUser.LastAppear.Add(60)
    err := ctx.UpdateLastAppear(userCopy.UserId, userCopy.LastAppear)
    if err != nil {
        t.Errorf("update err: %v", err)
    }

    updated, err := ctx.FindUser(userCopy.UserId)
    if err != nil {
        t.Errorf("update err: %v", err)
    }
    if updated.LastAppear.Before(userCopy.LastAppear) {
        t.Errorf("update err: time update failed")
    }

    // update
    userCopy.LastAppear = testUser.LastAppear.Add(-1)
    err = ctx.UpdateLastAppear(userCopy.UserId, userCopy.LastAppear)
    if err != nil {
        t.Errorf("update err: %v", err)
    }

    updated, err = ctx.FindUser(userCopy.UserId)
    if err != nil {
        t.Errorf("update err: %v", err)
    }
    if updated.LastAppear.Before(userCopy.LastAppear) {
        t.Errorf("update err: time update failed")
    }

    // update unknown user
    err = ctx.UpdateLastAppear("unknown userid", userCopy.LastAppear)
    if err == nil {
        t.Errorf("update err: unknown userid")
    }
}

func TestErase(t *testing.T) {

    // erase
    err := ctx.EraseUser(testUser.UserId)
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
