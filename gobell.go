package main

import (
    "./lease"
    "./line"
    "./user"
    "fmt"
    "regexp"
    "sync"
    "github.com/line/line-bot-sdk-go/linebot"
    "github.com/go-fsnotify/fsnotify"
    "github.com/BurntSushi/toml"
)

type Config struct {
    LeasePath string
}

var config Config

var currentUsers []user.User


func main() {

    _, err := toml.DecodeFile("config.toml", &config)
    if err != nil {
        fmt.Println(err)
        return
    }

    wg := &sync.WaitGroup{}
    wg.Add(1)

    // start file watching
    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        fmt.Println(err)
        return
    }
    defer watcher.Close()

    err = watcher.Add(config.LeasePath)
    if err != nil {
        fmt.Println(err)
        return
    }

    go func() {
        fmt.Println("start file watcher for:", config.LeasePath)
        for {
            select {
            case event := <-watcher.Events:
                fmt.Println("event:", event)
                watchEventHandler(event.Op, event.Name)
            case err := <-watcher.Errors:
                fmt.Println("watcher error: ", err)
                wg.Done()
            }
        }
    }()

    // bot server
    go func() {
        line.StartLineBotServer(lineEventHandler)
        wg.Done()
    }()

    wg.Wait()
}

func watchEventHandler(op fsnotify.Op, filename string) {

    if op&fsnotify.Remove == fsnotify.Remove {
        return
    }

    ctx := user.GetContext()
    defer ctx.Close()

    fmt.Printf("%v is modified!", filename)

    // load lease file
    leases, err := lease.Parse(config.LeasePath)
    if err != nil {
        fmt.Println(err)
        return
    }

    // update last appear time
    latestUsers := []user.User{}
    for _, l := range leases {
        user, err := ctx.FindMac(l.Mac)
        if err != nil {
            continue
        }
        ctx.UpdateLastAppear(user.UserId, *l.Start)
        latestUsers = append(latestUsers, *user)
    }
    defer func() {
        currentUsers = latestUsers
    }()

    cameUsers := []user.User{}
    for _, u := range latestUsers {
        if !contains(currentUsers, u.UserId) {
            cameUsers = append(cameUsers, u)
        }
    }

    leftUsers := []user.User{}
    for _, u := range cameUsers {
        if !contains(latestUsers, u.UserId) {
            leftUsers = append(leftUsers, u)
        }
    }

    allUserId, err := ctx.AllUserId()
    if err != nil {
        fmt.Println(err)
        return
    }
    bot := line.GetBotClient()

    // notify came members for all
    if len(cameUsers) > 0 {
        cameMes := ""
        for _, u := range cameUsers {
            cameMes += fmt.Sprintf("%vさん\n", u.Name)
        }
        cameMes += "が来ました"

        for _, userId := range allUserId {
            if _, err := bot.PushMessage(userId, linebot.NewTextMessage(cameMes)).Do(); err != nil{
                fmt.Println(err)
            }
        }
    }

    // notify left members for all
    if len(leftUsers) > 0 {
        leftMes := ""
        for _, u := range cameUsers {
            leftMes += fmt.Sprintf("%vさん\n", u.Name)
        }
        leftMes += "がいなくなりました"

        for _, userId := range allUserId {
            if _, err := bot.PushMessage(userId, linebot.NewTextMessage(leftMes)).Do(); err != nil{
                fmt.Println(err)
            }
        }
    }
}

func contains(users []user.User, userId string) bool {
    for _, u := range users {
        if u.UserId == userId {
            return true
        }
    }
    return false
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
            leases, err := lease.Parse(config.LeasePath)
            if err != nil {
                fmt.Println(err)
                return
            }

            var text = ""
            for _, u := range currentUsers {
                text += fmt.Sprintf("%v (%v)\n", u.Name, u.LastAppear)
            }
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
        if user == nil {
            continue
        }

        ret += fmt.Sprintln(user.Name)
    }

    return ret
}
