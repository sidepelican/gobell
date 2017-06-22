package line

import (
    "sort"
    "log"
    "path"
    "path/filepath"
    "time"
    "container/list"

    "github.com/sidepelican/gobell/config"
    "github.com/sidepelican/gobell/udb"
    "github.com/sidepelican/gobell/lease"

    "github.com/go-fsnotify/fsnotify"
)

var currentUsers udb.Users

var currentLogs *list.List = list.New()

func CurrentLogs() []WatchLog {
    ret := make([]WatchLog, currentLogs.Len())
    i := 0
    for e := currentLogs.Front(); e != nil; e = e.Next() {
        ret[i] = e.Value.(WatchLog)
        i ++
    }
    return ret
}

func StartFileWatcher() error {

    // watcher setup
    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        return err
    }
    defer watcher.Close()

    watchPath, _ := path.Split(config.LeasePath())
    err = watcher.Add(watchPath)
    if err != nil {
        return err
    }

    absPath, err := filepath.Abs(watchPath)
    if err != nil {
        absPath = watchPath
    }

    log.Println("start file watcher for:", absPath)
    for {
        select {
        case event := <-watcher.Events:
            watchEventHandler(event.Op, event.Name)
        case err := <-watcher.Errors:
            log.Println("watcher error:", err)
        }
    }

    return nil
}

func watchEventHandler(op fsnotify.Op, filePath string) {

    if filePath != config.LeasePath() {
        return
    }
    if op&fsnotify.Remove == fsnotify.Remove {
        return
    }

    // wait a minute to sum sequentially events
    if len(currentUsers) > 0 {
        time.Sleep(1 * time.Minute)
    }

    ctx := udb.GetContext()
    defer ctx.Close()

    // load lease file
    leases, err := lease.Parse(config.LeasePath())
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
            unregisteredUser.LastAppear = l.Start.Local()
            latestUsers = append(latestUsers, unregisteredUser)
            continue
        }
        ctx.UpdateLastAppear(u.UserId, *l.Start)
        latestUsers = append(latestUsers, *u)
    }
    sort.Sort(latestUsers)

    cameUsers := latestUsers.Difference(currentUsers)
    leftUsers := currentUsers.Difference(latestUsers)

    currentUsers = latestUsers

    // queue
    for _, u := range cameUsers {
        currentLogs.PushBack(WatchLog{u, CAME})
    }
    for _, u := range leftUsers {
        currentLogs.PushBack(WatchLog{u, LEFT})
    }
    for currentLogs.Len() > 100 {
        currentLogs.Remove(currentLogs.Front())
    }

    // notify
    NotifyCameAndLeftUsers(cameUsers, leftUsers)
}
