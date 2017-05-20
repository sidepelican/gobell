package line

import (
    "net/http"
    "log"
    "fmt"
    "sync"
    "regexp"
    "strings"

    "github.com/sidepelican/gobell/udb"
    "github.com/sidepelican/gobell/lease"

    "github.com/line/line-bot-sdk-go/linebot"
)

func HttpHandler(w http.ResponseWriter, r *http.Request) {

    bot := GetBotClient()

    // parse request
    events, err := bot.ParseRequest(r)
    if err != nil {
        // Do something when something bad happened.
        log.Println(err)
        return
    }

    wg := &sync.WaitGroup{}
    wg.Add(len(events))

    // handle event
    for _, event := range events {
        go func(bot *linebot.Client, event *linebot.Event) {
            defer wg.Done()
            lineEventHandler(bot, event)
        }(bot, event)
    }
    wg.Wait()

    fmt.Fprint(w, "{}")
}

func lineEventHandler(bot *linebot.Client, event *linebot.Event) {

    userId := event.Source.UserID

    switch event.Type {
    case linebot.EventTypeJoin:
        mes := linebot.NewTextMessage("ようこそ！最初にMacアドレスの登録をお願いしております。Macアドレスを入力してください✒️")
        _, err := bot.ReplyMessage(event.ReplyToken, mes).Do()
        if err != nil {
            log.Println(err)
            return
        }

    case linebot.EventTypeLeave:
        ctx := udb.GetContext()
        defer ctx.Close()

        err := ctx.EraseUser(userId)
        if err != nil {
            log.Println(err)
            return
        }

    case linebot.EventTypeMessage:

        replyText := func() string {
            ctx := udb.GetContext()
            defer ctx.Close()

            // for registered user
            _, err := ctx.FindUser(userId)
            if err == nil {
                // reply currentUsers

                if len(udb.CurrentUsers) == 0 {
                    return "誰もいないか、何かがおかしいようです"
                }

                const layout = "15:04"
                var text = ""
                for _, u := range udb.CurrentUsers {
                    text += fmt.Sprintf("%v (%v)\n", u.Name, u.LastAppear.Format(layout))
                }
                text = strings.TrimRight(text, "\n")
                return text
            }

            // for new visiter
            textMessage, ok := event.Message.(*linebot.TextMessage)
            if ok == false {
                // ignore photo, video, etc messages
                return ""
            }

            // check the message is mac addr or not
            macAddr, ok := findMacAddr(textMessage.Text)
            if ok == false {
                // mac addr not found. err reply
                return "Macアドレスを入力してください"
            }

            // check the mac addr is not registered
            _, err = ctx.FindMac(macAddr)
            if err == nil {
                return fmt.Sprintf("%vはすでに登録されています", macAddr)
            }

            // register

            // request username
            res, err := bot.GetProfile(userId).Do()
            if err != nil {
                log.Println(err)
                return err.Error()
            }

            // insert new user
            err = ctx.InsertUser(udb.NewUser(userId, macAddr, res.DisplayName))
            if err != nil {
                log.Println(err)
                return err.Error()
            }

            // insert succeeded
            return fmt.Sprintf("%vさんの登録が完了しました", res.DisplayName)
        }()

        if replyText != "" {
            _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyText)).Do()
            if err != nil {
                log.Println(err)
                return
            }
        }
    }
}

func NotifyCameAndLeftUsers(cameUsers udb.Users, leftUsers udb.Users) error {

    ctx := udb.GetContext()
    defer ctx.Close()

    allUserId, err := ctx.AllUserId()
    if err != nil {
        return err
    }
    bot := GetBotClient()

    // notify came members for all
    const layout = "15:04"
    if len(cameUsers) > 0 {
        cameMes := ""
        for _, u := range cameUsers {
            cameMes += fmt.Sprintf("%vさん(%v)\n", u.Name, u.LastAppear.Format(layout))
        }
        cameMes += "が来ました"

        for _, userId := range allUserId {
            if _, err := bot.PushMessage(userId, linebot.NewTextMessage(cameMes)).Do(); err != nil {
                log.Println(err)
            }
        }
    }

    // notify left members for all
    if len(leftUsers) > 0 {
        leftMes := ""
        for _, u := range leftUsers {
            leftMes += fmt.Sprintf("%vさん(%v)\n", u.Name, u.LastAppear.Format(layout))
        }
        leftMes += "がいなくなりました"

        for _, userId := range allUserId {
            if _, err := bot.PushMessage(userId, linebot.NewTextMessage(leftMes)).Do(); err != nil {
                log.Println(err)
            }
        }
    }

    return nil
}

func findMacAddr(message string) (string, bool) {

    macRegex := regexp.MustCompile(`^([0-9A-Fa-f]{2}[:-]){5}[0-9A-Fa-f]{2}$`)
    macResult := macRegex.FindString(message)

    if macResult == "" {
        return "", false
    }

    macResult = lease.TrimMacAddr(macResult)
    return macResult, true
}
