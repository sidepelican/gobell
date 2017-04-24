package main

import (
    "testing"

    "log"
    "os"
    "fmt"

    "github.com/sidepelican/gobell/lease"
    "github.com/sidepelican/gobell/udb"

    "github.com/BurntSushi/toml"
)

func setup() {
    println("setup")

    _, err := toml.DecodeFile(getRunPath() + "config.toml", &config)
    if err != nil {
        log.Println(err)
        return
    }
}

func teardown() {
    println("teardown")

}

func TestUnregistered(t *testing.T) {

    // load lease file
    leases, err := lease.Parse(config.LeasePath)
    if err != nil {
        log.Println(err)
        return
    }

    // update last appear time
    latestUsers := udb.Users{}
    for _, l := range leases {

        // unregistered user
        unregisteredUser := udb.NewUser(l.Mac, l.Mac, l.Hostname)
        unregisteredUser.LastAppear = *l.Start
        latestUsers = append(latestUsers, unregisteredUser)
    }
    currentUsers = latestUsers

    const layout = "15:04"
    for _, u := range currentUsers {
        fmt.Printf("%v (%v)\n", u.Name, u.LastAppear.Format(layout))
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