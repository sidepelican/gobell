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

    // if already have a token, update it and return.
    if authToken := r.Header.Get("Authorization"); authToken != "" {
        newToken, err := auth.UpdateExpires(authToken)
        if err == nil {
            redererer.JSON(w, http.StatusOK, CommonResponse{
                Status:  http.StatusOK,
                Message: newToken,
            })
            return
        }
    }

    name := r.FormValue("name")
    pass := r.FormValue("pass")

    token, err := auth.Auth(name, pass)
    if err != nil {
        message := "Authorization failed! Reason: " + err.Error()
        redererer.JSON(w, http.StatusNotAcceptable, NewErrorResponse(http.StatusNotAcceptable, message))
        return
    }

    redererer.JSON(w, http.StatusOK, CommonResponse{
        Status:  http.StatusOK,
        Message: token,
    })
}
