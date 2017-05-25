package handler

import (
    "context"
    "net/http"
    "time"

    "github.com/unrolled/render"
)

const TimeOutLimit = 30 * time.Second
var redererer = render.New(render.Options{IndentJSON: true})

type ErrorResponse struct {
    Status  int    `json:"status"`
    Message string `json:"message"`
}

func NewErrorResponse(status int, message string) ErrorResponse {
    return ErrorResponse{
        Status:  status,
        Message: message,
    }
}

type HandlerWithDI struct {
    Impl func(http.ResponseWriter, *http.Request)
}

func (h HandlerWithDI) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(r.Context(), TimeOutLimit)
    defer cancel()
    cr := r.WithContext(ctx)
    h.Impl(w, cr)
}