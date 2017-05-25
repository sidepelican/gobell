package main

import (
    "sync"
    "log"
    "net/http"

    "github.com/sidepelican/gobell/line"
    "github.com/sidepelican/gobell/config"
    "github.com/sidepelican/gobell/handler"

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
        if err := line.StartFileWatcher(); err != nil {
            log.Println("Watcher:", err)
        }
        wg.Done()
    }()

    // bot server
    go func() {
        r := mux.NewRouter()
        r.HandleFunc("/line", line.HttpHandler)
        r.HandleFunc("/list", handler.ListHandler)
        user := r.PathPrefix("/user").Subrouter()
        user.HandleFunc("/list", handler.UserListHandler)
        user.HandleFunc("/add", handler.UserAddHandler)

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
