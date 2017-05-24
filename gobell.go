package main

import (
    "sync"
    "log"
    "net/http"
    "encoding/json"
    "sort"

    "github.com/sidepelican/gobell/line"
    "github.com/sidepelican/gobell/config"
    "github.com/sidepelican/gobell/watch"
    "github.com/sidepelican/gobell/lease"
    "github.com/sidepelican/gobell/udb"

    "github.com/gorilla/mux"
)

func main() {

    log.SetFlags(log.Lshortfile)

    if err := config.InitConfig(); err != nil {
        log.Println(err)
        return
    }

    wg := &sync.WaitGroup{}
    wg.Add(1)

    // start file watching
    go func() {
        if err := watch.StartFileWatcher(); err != nil {
            log.Println("Watcher:", err)
        }
        wg.Done()
    }()

    // bot server
    go func() {
        r := mux.NewRouter()
        r.HandleFunc("/line", line.HttpHandler)
        r.HandleFunc("/list", listHandler)
        user := r.PathPrefix("/user").Subrouter()
        user.HandleFunc("/list", userListHandler)

        srv := &http.Server{
            Addr:    ":8080",
            Handler: r,
        }

        log.Println("starting linebot server")
        if err := srv.ListenAndServe(); err != nil {
            log.Println("ListenAndServe: ", err)
        }
        wg.Done()
    }()

    wg.Wait()
}

func listHandler(w http.ResponseWriter, r *http.Request) {

    // load lease file
    leases, err := lease.Parse(config.LeasePath())
    if err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    sort.Sort(leases)

    bytes, err := json.Marshal(leases)
    if err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write(bytes)
}

func userListHandler(w http.ResponseWriter, r *http.Request) {

    ctx := udb.GetContext()

    users, err := ctx.AllUsers()
    if err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    if len(users) == 0 {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("{}"))
        return
    }

    sort.Sort(users)

    bytes, err := json.Marshal(users)
    if err != nil {
        log.Println(err)
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write(bytes)
}
