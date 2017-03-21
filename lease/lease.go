package lease

import (
    "fmt"
    "time"
)

type Lease struct {
    start    *time.Time
    ip       string
    mac      string
    hostname string
}

type LeaseFinder struct {
    results map[string]Lease // key: ip Addr
    current *Lease
}

const layout = "2006/01/02 15:04:05"

func NewLeaseFinder() LeaseFinder {
    return LeaseFinder{
        results: make(map[string]Lease),
        current: nil,
    }
}

func (f *LeaseFinder)FindStart(ip string) {

    if f.current != nil {
        fmt.Println("findStartLease called before call findEnd. something wrong")
        f.current = nil
    }

    f.current = &Lease{ip: ip}
}

func (f *LeaseFinder)FindStartTime(startString string) {

    t, err := time.Parse(layout, startString)
    if err != nil {
        fmt.Println(err)
        return
    }

    f.current.start = &t
}

func (f *LeaseFinder)FindMac(mac string) {
    f.current.mac = mac
}

func (f *LeaseFinder)FindHostname(hostname string) {
    f.current.hostname = hostname
}

func (f *LeaseFinder)FindEnd() {

    defer func() {
        f.current = nil
    }()

    // something dropped
    if f.current.hostname == "" || f.current.start == nil || f.current.ip == "" || f.current.mac == "" {
        return
    }

    // all values are completed

    old, ok := f.results[f.current.ip]
    if ok {
        if old.start.After(*f.current.start) {
            return
        }
    }

    f.results[f.current.ip] = *f.current

    return
}

func (f *LeaseFinder)PrintAll() {

    for _, v := range f.results {
        fmt.Println(v.ip + "," + v.start.Format(layout) + "," + v.mac + "," + v.hostname)
    }
}