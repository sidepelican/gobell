package main

import (
    "./lease"
    "./line"
    "./user"

    //"github.com/go-fsnotify/fsnotify"
    "github.com/line/line-bot-sdk-go/linebot"
    "fmt"
    "regexp"
    "sync"
)

func main() {
    wg := &sync.WaitGroup{}
    wg.Add(1)
    go func() {
        line.StartLineBotServer(lineEventHandler)
        wg.Done()
    }()
    wg.Wait()
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

    case linebot.EventTypeLeave:
        ctx := user.GetContext()
        defer ctx.Close()

        err := ctx.EraseUser(userId)
        if err != nil {
            fmt.Println(err)
            return
        }

    case linebot.EventTypeMessage:
        ctx := user.GetContext()
        defer ctx.Close()

        // for registered user
        _, err := ctx.FindUser(userId)
        if err == nil {

            // load dhcpd.lease
            leases, err := lease.Parse("sample/dhcpd.lease")
            if err != nil {
                fmt.Println(err)
                return
            }

            var text = ""
            text += fmt.Sprintln(registeredUserNames(leases))
            text += "-------------\n"
            text += leases.AllHostname()

            message := linebot.NewTextMessage(text)
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
            // ignore photo, video, etc messages
            return
        }

        // check the message is mac addr or not
        macAddr, ok := isMacAddr(textMessage.Text)
        if ok == false {
            // mac addr not found. err reply
            message := linebot.NewTextMessage("Macアドレスを入力してください")
            _, err := bot.ReplyMessage(event.ReplyToken, message).Do()
            if err != nil {
                fmt.Println(err)
                return
            }
            return
        }

        // check the mac addr is not registered
        _, err = ctx.FindMac(macAddr)
        if err == nil {
            message := linebot.NewTextMessage(fmt.Sprintf("%vはすでに登録されています", macAddr))
            _, err := bot.ReplyMessage(event.ReplyToken, message).Do()
            if err != nil {
                fmt.Println(err)
                return
            }
            return
        }

        // register

        // request username
        res, err := bot.GetProfile(userId).Do()
        if err != nil {
            fmt.Println(err)
            return
        }

        // insert new user
        err = ctx.InsertUser(user.NewUser(userId, macAddr, res.DisplayName))
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
    }
}

func isMacAddr(message string) (string, bool) {

    macRegex := regexp.MustCompile(`^([0-9A-Fa-f]{2}[:-]){5}[0-9A-Fa-f]{2}$`)
    macResult := macRegex.FindString(message)

    if macResult == "" {
        return "", false
    }

    return macResult, true
}

func registeredUserNames(leases lease.Leases) string {

    ctx := user.GetContext()
    defer ctx.Close()

    var ret = ""
    for _, l := range leases {
        user, _ := ctx.FindMac(l.Mac)
        if user == nil { continue }

        ret += fmt.Sprintln(user.Name)
    }

    return ret
}