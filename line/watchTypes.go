package line

import (
    "github.com/sidepelican/gobell/udb"
    "fmt"
)

type LogType int

const (
    CAME LogType = iota
    LEFT
)

type WatchLog struct {
    user    udb.User
    logType LogType
}

func (l *WatchLog) String() string {
    const layout = "15:04"
    u := l.user
    return fmt.Sprintf("%vさん(%v)が%v", u.Name, u.LastAppear.Format(layout), l.logType)
}

func (t LogType) String() string {
    switch t {
    case CAME:
        return "来ました"
    case LEFT:
        return "いなくなりました"
    default:
        return "???"
    }
}
