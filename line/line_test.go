package line

import (
    "github.com/line/line-bot-sdk-go/linebot"
    "testing"
)

const testUId = "<tester userID>"

func TestPushMessage(t *testing.T) {

    bot := GetBotClient()

    mes := linebot.NewTextMessage("push message test")
    _, err := bot.PushMessage(testUId, mes).Do()
    if err != nil {
        t.Error(err)
    }
}