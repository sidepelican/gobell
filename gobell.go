package main

import (
    "sync"
    "log"
    "github.com/sidepelican/gobell/line"
    "github.com/sidepelican/gobell/config"
    "github.com/sidepelican/gobell/watch"

    "github.com/gorilla/mux"
    "net/http"
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
