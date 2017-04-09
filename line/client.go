package line

import (
    "fmt"
    "sync"
    "github.com/line/line-bot-sdk-go/linebot"
    "github.com/BurntSushi/toml"
)

type LineConfig struct {
    Secret string
    Token  string
}

var once sync.Once
var bot *linebot.Client

func GetBotClient() *linebot.Client {

    once.Do(func() {
        var config LineConfig
        _, err := toml.DecodeFile("line/lineserver.toml", &config)
        if err != nil {
            fmt.Println(err)
            return
        }

        // create bot
        bot, err = linebot.New(config.Secret, config.Token)
        if err != nil {
            fmt.Println(err)
            return
        }
    })
    return bot
}
