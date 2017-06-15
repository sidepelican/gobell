package handler

import (
    "net/http"
    "github.com/sidepelican/gobell/auth"
)

type AuthHandler struct {
    Impl func(http.ResponseWriter, *http.Request)
}

func (h *AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

    authToken := r.Header.Get("Authorization")

    _, err := auth.Validate(authToken)
    if err != nil {
        redererer.JSON(w, http.StatusUnauthorized, NewErrorResponse(http.StatusUnauthorized, "Authorization required!"))
        return
    }

    h.Impl(w, r)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {

    name := r.FormValue("name")
    pass := r.FormValue("pass")

    token, err := auth.Auth(name, pass)
    if err != nil {
        redererer.JSON(w, http.StatusNotAcceptable, NewErrorResponse(http.StatusNotAcceptable, "Authorization failed!"))
        return
    }

    redererer.JSON(w, http.StatusOK, CommonResponse{
        Status:  http.StatusOK,
        Message: token,
    })
}