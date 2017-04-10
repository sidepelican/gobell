package main

import (
    "./lease"
    "./line"
    "./user"

    //"github.com/go-fsnotify/fsnotify"
    "github.com/line/line-bot-sdk-go/linebot"
    "fmt"
    "regexp"
)

func main() {
    line.StartLineBotServer(lineEventHandler)
}

func lineEventHandler(bot *linebot.Client, event *linebot.Event) {

    userId := event.Source.UserID

    switch event.Type {
    case linebot.EventTypeJoin:
        mes := linebot.NewTextMessage("ようこそ！最初にMacアドレスの登録をお願いしております。Macアドレスを入力してください✒️")
        _, err := bot.ReplyMessage(event.ReplyToken, mes).Do()
        if err != nil {
            fmt.Println(err)
            return
        }

    case linebot.EventTypeMessage:

        ctx := user.GetContext()
        defer ctx.Close()
        currentUser, _ := ctx.FindUser(userId)

        // for registered user
        if currentUser != nil {

            // load dhcpd.lease
            leases, err := lease.Parse("sample/dhcpd.lease")
            if err != nil {
                fmt.Println(err)
                return
            }
            hostnames := lease.AllHostname(leases)

            message := linebot.NewTextMessage(hostnames)
            _, err = bot.ReplyMessage(event.ReplyToken, message).Do()
            if err != nil {
                fmt.Println(err)
                return
            }

            return
        }

        // for new visiter
        textMessage, ok := event.Message.(*linebot.TextMessage)
        if ok == false {
            // does not reaction
            return
        }

        // check the message is mac addr or not
        macRegex := regexp.MustCompile(`^([0-9A-Fa-f]{2}[:-]){5}[0-9A-Fa-f]{2}$`)
        macResult := macRegex.FindString(textMessage.Text)

        if macResult != "" {
            // register

            // fetch username
            res, err := bot.GetProfile(userId).Do()
            if err != nil {
                fmt.Println(err)
                return
            }

            // check the mac addr is not registered
            _, err = ctx.FindMac(macResult)
            if err == nil {
                message := linebot.NewTextMessage(fmt.Sprintf("%vはすでに登録されています", macResult))
                _, err := bot.ReplyMessage(event.ReplyToken, message).Do()
                if err != nil {
                    fmt.Println(err)
                    return
                }
                return
            }

            // insert new user
            err = ctx.InsertUser(user.NewUser(userId, macResult, res.DisplayName))
            if err != nil {
                fmt.Println(err)
                return
            }

            // insert succeeded
            message := linebot.NewTextMessage(fmt.Sprintf("%vさんの登録が完了しました", res.DisplayName))
            _, err = bot.ReplyMessage(event.ReplyToken, message).Do()
            if err != nil {
                fmt.Println(err)
                return
            }
            return

        }else{
            // mac addr not found. err reply
            message := linebot.NewTextMessage("Macアドレスを入力してください")
            _, err := bot.ReplyMessage(event.ReplyToken, message).Do()
            if err != nil {
                fmt.Println(err)
                return
            }
            return
        }
    }
}
