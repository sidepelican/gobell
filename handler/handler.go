package handler

import (
    "log"
    "sort"
    "net/http"

    "github.com/sidepelican/gobell/config"
    "github.com/sidepelican/gobell/lease"
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
