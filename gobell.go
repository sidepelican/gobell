package main

import (
    "./lease"
    "./line"

    //"github.com/go-fsnotify/fsnotify"
    "github.com/line/line-bot-sdk-go/linebot"
    "fmt"
)

func main() {
    line.StartLineBotServer(lineEventHandler)
}

func lineEventHandler(bot *linebot.Client, event *linebot.Event) {

    // load dhcpd.lease
    leases, err := lease.Parse("sample/dhcpd.lease")
    if err != nil {
        fmt.Println(err)
        return
    }
    hostnames := lease.AllHostname(leases)
    message := linebot.NewTextMessage(hostnames)

    switch event.Type {
    case linebot.EventTypeJoin:
        mes := linebot.NewTextMessage("ようこそ！最初にMacアドレスの登録をお願いしております。Macアドレスを入力してください✒️")
        _, err := bot.ReplyMessage(event.ReplyToken, mes).Do()
        if err != nil {
            fmt.Println(err)
            return
        }

    case linebot.EventTypeMessage:
        _, err := bot.ReplyMessage(event.ReplyToken, message).Do()
        if err != nil {
            fmt.Println(err)
            return
        }
    }
}