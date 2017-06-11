package main

import (
    "sync"
    "log"
    "net/http"

    "github.com/sidepelican/gobell/line"
    "github.com/sidepelican/gobell/config"
    "github.com/sidepelican/gobell/handler"

    "github.com/gorilla/mux"
    "github.com/rs/cors"
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

    // api & bot server
    go func() {
        r := mux.NewRouter()
        r.HandleFunc("/line", line.HttpHandler)
        r.HandleFunc("/list", handler.ListHandler)
        r.HandleFunc("/users", handler.UsersHandler)
        user := r.PathPrefix("/user").Subrouter()
        user.HandleFunc("/list", handler.UserListHandler)
        user.HandleFunc("/add", handler.UserAddHandler).Methods(http.MethodPost)
        user.HandleFunc("/delete", handler.UserDeleteHandler).Methods(http.MethodPost)
        user.HandleFunc("/note", handler.EditNoteHandler).Methods(http.MethodPost)

        srv := &http.Server{
            Addr:    ":8080",
            Handler: cors.Default().Handler(r),
        }

        if err := srv.ListenAndServe(); err != nil {
            log.Println("ListenAndServe: ", err)
        }
        wg.Done()
    }()

    wg.Wait()
}
