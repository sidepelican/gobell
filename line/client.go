package line

import (
    "sync"
    "github.com/line/line-bot-sdk-go/linebot"
    "log"
    "os"
)

var once sync.Once
var bot *linebot.Client

func GetBotClient() *linebot.Client {

    once.Do(func() {

        // create bot
        lineToken := os.Getenv("LINE_BOT_TOKEN")
        lineSecret := os.Getenv("LINE_BOT_SECRET")

        var err error
        bot, err = linebot.New(lineSecret, lineToken)
        if err != nil {
            log.Fatal(err)
            return
        }
    })
    return bot
}
