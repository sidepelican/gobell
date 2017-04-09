package line

import (
    "net/http"
    "log"
    "fmt"

    "../lease"

    "github.com/line/line-bot-sdk-go/linebot"
)

func StartLineBotServer() {
    http.HandleFunc("/", handler)
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal("ListenAndServe: ", err)
        return
    }
}

func handler(w http.ResponseWriter, r *http.Request) {

    bot := GetBotClient()

    // parse request
    events, err := bot.ParseRequest(r)
    if err != nil {
        // Do something when something bad happened.

        fmt.Println(err)
        return
    }

    // load dhcpd.lease
    leases, err := lease.Parse("sample/dhcpd.lease")
    if err != nil {
        fmt.Println(err)
        return
    }
    hostnames := lease.AllHostname(leases)
    message := linebot.NewTextMessage(hostnames)

    // handle event
    for _, event := range events {
        if event.Type == linebot.EventTypeMessage {

            replyToken := event.ReplyToken
            _, err = bot.ReplyMessage(replyToken, message).Do()
            if err != nil {
                fmt.Println(err)
                break
            }
        }
    }

    fmt.Fprintf(w, "")
}