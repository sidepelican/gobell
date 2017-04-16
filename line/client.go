package line

import (
    "sync"
    "github.com/line/line-bot-sdk-go/linebot"
    "github.com/BurntSushi/toml"
    "log"
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
            log.Println(err)
            return
        }

        // create bot
        bot, err = linebot.New(config.Secret, config.Token)
        if err != nil {
            log.Println(err)
            return
        }
    })
    return bot
}
