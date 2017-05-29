package lease

import (
    "fmt"
    "os"
    "regexp"
    "bufio"
    "io"
    "time"
    "log"
)

const layout = "2006/01/02 15:04:05"
var timeZone = time.UTC

type leaseFinder struct {
    results map[string]Lease // key: ip Addr
    current *Lease
}

func Parse(path string) (Leases, error) {

    // file open
    fp, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer fp.Close()

    // regexp setup
    rLease := regexp.MustCompile(`lease ([0-9¥.]+) {`)
    rStarts := regexp.MustCompile(`starts ([0-9]) (.*);`)
    rHwEth := regexp.MustCompile(`hardware ethernet ([0-9A-Fa-f:-]+);`)
    rHostname := regexp.MustCompile(`client-hostname "(.*)";`)

    lf := newLeaseFinder()

    // start reading
    reader := bufio.NewReaderSize(fp, 4096)
    for {
        lineBuf, _, err := reader.ReadLine()
        if err == io.EOF {
            break
        } else if err != nil {
            return nil, err
        }

        line := string(lineBuf)

        // try each regexp
        var res []string = nil
        res = rLease.FindStringSubmatch(line)
        if res != nil {
            lf.FindStart(res[len(res)-1])
        }
        res = rStarts.FindStringSubmatch(line)
        if res != nil {
            lf.FindStartTime(res[len(res)-1])
        }
        res = rHwEth.FindStringSubmatch(line)
        if res != nil {
            lf.FindMac(res[len(res)-1])
        }
        res = rHostname.FindStringSubmatch(line)
        if res != nil {
            lf.FindHostname(res[len(res)-1])
        }
        if line == `}` {
            lf.FindEnd()
        }
    }

    return lf.AllLeases(), nil
}

func newLeaseFinder() leaseFinder {
    return leaseFinder{
        results: make(map[string]Lease),
        current: nil,
    }
}

func (f *leaseFinder)FindStart(ip string) {

    if f.current != nil {
        log.Println("findStartLease called before call findEnd. something wrong")
        f.current = nil
    }

    f.current = &Lease{Ip: ip}
}

func (f *leaseFinder)FindStartTime(startString string) {

    t, err := time.ParseInLocation(layout, startString, timeZone)
    if err != nil {
        log.Println(err)
        return
    }

    localTime := t.Local()
    f.current.Start = &localTime
}

func (f *leaseFinder)FindMac(mac string) {
    f.current.Mac = TrimMacAddr(mac)
}

func (f *leaseFinder)FindHostname(hostname string) {
    f.current.Hostname = hostname
}

func (f *leaseFinder)FindEnd() {

    defer func() {
        f.current = nil
    }()

    // something dropped
    if f.current.Hostname == "" || f.current.Start == nil || f.current.Ip == "" || f.current.Mac == "" {
        return
    }

    // all values are completed
    old, ok := f.results[f.current.Ip]
    if ok {
        // when the "start time" is not latest, will be ignored
        if old.Start.After(*f.current.Start) {
            return
        }
    }

    f.results[f.current.Ip] = *f.current

    return
}

func (f *leaseFinder)PrintAll() {
    for _, v := range f.results {
        fmt.Println(v.Ip + "," + v.Start.Format(layout) + "," + v.Mac + "," + v.Hostname)
    }
}

func (f *leaseFinder)AllLeases() Leases {
    ret := make([]Lease, 0, len(f.results))
    for _, v := range f.results {
        ret = append(ret, v)
    }
    return ret
}