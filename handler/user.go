package handler

import (
    "sort"
    "log"
    "net/http"

    "github.com/sidepelican/gobell/udb"
    "github.com/sidepelican/gobell/lease"
)

func UserListHandler(w http.ResponseWriter, r *http.Request) {
    ctx := udb.GetContext()
    defer ctx.Close()

    users, err := ctx.AllUsers()
    if err != nil {
        log.Println(err)
        redererer.JSON(w, http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, err.Error()))
        return
    }

    if len(users) == 0 {
        redererer.Text(w, http.StatusOK, "{}")
        return
    }

    sort.Sort(users)
    redererer.JSON(w, http.StatusOK, users)
}

func UserAddHandler(w http.ResponseWriter, r *http.Request) {
    ctx := udb.GetContext()
    defer ctx.Close()

    name := r.FormValue("name")
    mac := lease.TrimMacAddr(r.FormValue("mac"))
    note := r.FormValue("note")

    if name == "" || mac == "" {
        mes := "name or mac address incorrect."
        redererer.JSON(w, http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, mes))
        return
    }

    user := udb.NewUser(mac, mac, name)
    user.Note = note
    err := ctx.InsertUser(user)
    if err != nil {
        log.Println(err)
        redererer.JSON(w, http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, err.Error()))
        return
    }

    redererer.JSON(w, http.StatusOK, NewSuccessResponse())
}

func UserDeleteHandler(w http.ResponseWriter, r *http.Request) {
    ctx := udb.GetContext()
    defer ctx.Close()

    err := r.ParseForm()
    if err != nil {
        mes := "failed to parse form: " + err.Error()
        redererer.JSON(w, http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, mes))
        return
    }

    ids := r.PostForm["user_ids[]"]
    if len(ids) == 0 {
        mes := "userId not found."
        redererer.JSON(w, http.StatusBadRequest, NewErrorResponse(http.StatusBadRequest, mes))
        return
    }

    for _, userId := range ids {
        err := ctx.EraseUser(userId)
        if err != nil {
            log.Println(err)
            redererer.JSON(w, http.StatusInternalServerError, NewErrorResponse(http.StatusInternalServerError, err.Error()))
            return
        }
    }

    redererer.JSON(w, http.StatusOK, NewSuccessResponse())
}
