package line

import (
    "net/http"
    "log"
    "fmt"

    "github.com/line/line-bot-sdk-go/linebot"
    "sync"
)

var linebotHandler func(*linebot.Client, *linebot.Event)

func StartLineBotServer(handler func(*linebot.Client, *linebot.Event)) {

    linebotHandler = handler

    fmt.Println("starting linebot server")
    http.HandleFunc("/", httpHandler)
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal("ListenAndServe: ", err)
        return
    }
}

func httpHandler(w http.ResponseWriter, r *http.Request) {

    bot := GetBotClient()

    // parse request
    events, err := bot.ParseRequest(r)
    if err != nil {
        // Do something when something bad happened.

        fmt.Println(err)
        return
    }

    wg := &sync.WaitGroup{}
    wg.Add(len(events))

    // handle event
    for _, event := range events {
        go func(bot *linebot.Client, event *linebot.Event) {
            defer wg.Done()
            linebotHandler(bot, event)
        }(bot, event)
    }
    wg.Wait()

    fmt.Fprint(w, "{}")
}