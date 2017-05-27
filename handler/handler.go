package handler

import (
    "log"
    "sort"
    "net/http"

    "github.com/sidepelican/gobell/config"
    "github.com/sidepelican/gobell/lease"
    "github.com/sidepelican/gobell/udb"
)

func ListHandler(w http.ResponseWriter, r *http.Request) {

    // load lease file
    leases, err := lease.Parse(config.LeasePath())
    if err != nil {
        log.Println(err)
        redererer.JSON(w, http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, err.Error()))
        return
    }

    if len(leases) == 0 {
        redererer.Text(w, http.StatusOK, "{}")
        return
    }

    sort.Sort(leases)
    redererer.JSON(w, http.StatusOK, leases)
}

func UsersHandler(w http.ResponseWriter, r *http.Request) {

    // load lease file
    leases, err := lease.Parse(config.LeasePath())
    if err != nil {
        log.Println(err)
        redererer.JSON(w, http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, err.Error()))
        return
    }

    if len(leases) == 0 {
        redererer.Text(w, http.StatusOK, "{}")
        return
    }

    ctx := udb.GetContext()
    defer ctx.Close()

    var currentUsers udb.Users
    for _, l := range leases {
        ctx.UpdateLastAppear(l.Mac, *l.Start)
        if u, _ := ctx.FindMac(l.Mac); u != nil {
            currentUsers = append(currentUsers, *u)
        }
    }
    currentUsers = removeDuplicate(currentUsers)

    sort.Sort(currentUsers)
    redererer.JSON(w, http.StatusOK, currentUsers)
}

func removeDuplicate(users udb.Users) udb.Users {
    len := len(users)
    results := make(udb.Users, 0, len)
    encountered := map[string]udb.User{}
    for i := 0; i < len; i++ {
        encount, ok := encountered[users[i].Name]
        if ok == false || encount.LastAppear.Before(users[i].LastAppear) {
            encountered[users[i].Name] = users[i]
            results = append(results, users[i])
        }
    }
    return results
}