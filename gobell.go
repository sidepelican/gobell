package main

import (
    "fmt"
    "regexp"
    "sync"
    "log"
    "path"
    "sort"
    "path/filepath"
    "os"
    "strings"

    "github.com/sidepelican/gobell/lease"
    "github.com/sidepelican/gobell/udb"
    "github.com/sidepelican/gobell/line"

    "github.com/line/line-bot-sdk-go/linebot"
    "github.com/go-fsnotify/fsnotify"
    "github.com/BurntSushi/toml"
)

type Config struct {
    LeasePath string
}
var config Config

var currentUsers udb.Users

func main() {

    log.SetFlags(log.Lshortfile | log.LstdFlags)

    _, err := toml.DecodeFile(getRunPath() + "config.toml", &config)
    if err != nil {
        log.Println(err)
        return
    }

    wg := &sync.WaitGroup{}
    wg.Add(1)

    // start file watching
    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        log.Println(err)
        return
    }
    defer watcher.Close()

    leaseDir, _ := path.Split(config.LeasePath)
    err = watcher.Add(leaseDir)
    if err != nil {
        log.Println(err)
        return
    }

    go func() {
        log.Println("start file watcher for:", config.LeasePath)
        for {
            select {
            case event := <-watcher.Events:
                log.Println("watcher:", event)
                watchEventHandler(event.Op, event.Name)
            case err := <-watcher.Errors:
                log.Println("watcher error:", err)
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

    if filename != config.LeasePath {
        return
    }
    if op&fsnotify.Remove == fsnotify.Remove {
        return
    }

    ctx := udb.GetContext()
    defer ctx.Close()

    // load lease file
    leases, err := lease.Parse(config.LeasePath)
    if err != nil {
        log.Println(err)
        return
    }

    // update last appear time
    latestUsers := udb.Users{}
    for _, l := range leases {
        u, _ := ctx.FindMac(l.Mac)
        if u == nil {
            // unregistered user
            unregisteredUser := udb.NewUser(l.Mac, l.Mac, l.Hostname)
            unregisteredUser.LastAppear = *l.Start
            latestUsers = append(latestUsers, unregisteredUser)
            continue
        }
        ctx.UpdateLastAppear(u.UserId, *l.Start)
        latestUsers = append(latestUsers, *u)
    }
    sort.Sort(latestUsers)

    cameUsers := udb.Users{}
    for _, u := range latestUsers {
        if !contains(currentUsers, u.UserId) {
            cameUsers = append(cameUsers, u)
        }
    }

    leftUsers := udb.Users{}
    for _, u := range currentUsers {
        if !contains(latestUsers, u.UserId) {
            leftUsers = append(leftUsers, u)
        }
    }

    currentUsers = latestUsers

    // notify
    allUserId, err := ctx.AllUserId()
    if err != nil {
        log.Println(err)
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
            if _, err := bot.PushMessage(userId, linebot.NewTextMessage(cameMes)).Do(); err != nil {
                log.Println(err)
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
            if _, err := bot.PushMessage(userId, linebot.NewTextMessage(leftMes)).Do(); err != nil {
                log.Println(err)
            }
        }
    }
}

func contains(users udb.Users, userId string) bool {
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
        ctx := udb.GetContext()
        defer ctx.Close()

        // for registered user
        _, err := ctx.FindUser(userId)
        if err == nil {
            // reply currentUsers
            const layout = "15:04"
            var text = ""
            for _, u := range currentUsers {
                text += fmt.Sprintf("%v (%v)\n", u.Name, u.LastAppear.Format(layout))
            }
            text = strings.TrimRight(text, "\n")

            message := linebot.NewTextMessage(text)
            _, err = bot.ReplyMessage(event.ReplyToken, message).Do()
            if err != nil {
                log.Println(err)
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
                log.Println(err)
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
                log.Println(err)
                return
            }
            return
        }

        // register

        // request username
        res, err := bot.GetProfile(userId).Do()
        if err != nil {
            log.Println(err)
            return
        }

        // insert new user
        err = ctx.InsertUser(udb.NewUser(userId, macAddr, res.DisplayName))
        if err != nil {
            log.Println(err)
            return
        }

        // insert succeeded
        message := linebot.NewTextMessage(fmt.Sprintf("%vさんの登録が完了しました", res.DisplayName))
        _, err = bot.ReplyMessage(event.ReplyToken, message).Do()
        if err != nil {
            log.Println(err)
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

    macResult = lease.TrimMacAddr(macResult)
    return macResult, true
}

func registeredUserNames(leases lease.Leases) string {

    ctx := udb.GetContext()
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

func getRunPath() string {
    dir, err := filepath.Abs(filepath.Dir(os.Args[0])) // Get the absolute path at Executing file. Reference：http://stackoverflow.com/questions/18537257/golang-how-to-get-the-directory-of-the-currently-running-file
    if err != nil {
        log.Println(err)
        return ""
    }

    // for `$go run ~~` support
    if strings.HasPrefix(dir, "/var") {
        return ""
    }

    return dir + "/"
}