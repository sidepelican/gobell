package lease

import (
    "os"
    "fmt"
    "time"
    "bufio"
    "io"
    "regexp"
    "strings"
)

type Lease struct {
    start    *time.Time
    ip       string
    mac      string
    hostname string
}

type leaseFinder struct {
    results map[string]Lease // key: ip Addr
    current *Lease
}

const layout = "2006/01/02 15:04:05"

// TODO: fat
func ParseLease(path string) []Lease {
    fp, err := os.Open(path)
    if err != nil {
        fmt.Println(err)
        return nil
    }
    defer fp.Close()

    // regexp setup
    rLease := regexp.MustCompile(`lease ([0-9¥.]+) {`)
    rStarts := regexp.MustCompile(`starts ([0-9]) (.*);`)
    rHwEth := regexp.MustCompile(`hardware ethernet ([0-9A-Fa-f:-]+);`)
    rHostname := regexp.MustCompile(`client-hostname "(.*)";`)
    rEnd := regexp.MustCompile(`}`)

    lf := NewLeaseFinder()

    // start reading
    reader := bufio.NewReaderSize(fp, 4096)
    for {
        lineBuf, _, err := reader.ReadLine()
        line := string(lineBuf)

        if err == io.EOF {
            break
        } else if err != nil {
            panic(err)
        }

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
        if rEnd.MatchString(line) {
            lf.FindEnd()
        }
    }

    return lf.AllLeases()
}

func NewLeaseFinder() leaseFinder {
    return leaseFinder{
        results: make(map[string]Lease),
        current: nil,
    }
}

func (f *leaseFinder)FindStart(ip string) {

    if f.current != nil {
        fmt.Println("findStartLease called before call findEnd. something wrong")
        f.current = nil
    }

    f.current = &Lease{ip: ip}
}

func (f *leaseFinder)FindStartTime(startString string) {

    t, err := time.Parse(layout, startString)
    if err != nil {
        fmt.Println(err)
        return
    }

    f.current.start = &t
}

func (f *leaseFinder)FindMac(mac string) {
    f.current.mac = mac
}

func (f *leaseFinder)FindHostname(hostname string) {
    f.current.hostname = hostname
}

func (f *leaseFinder)FindEnd() {

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

func (f *leaseFinder)PrintAll() {

    for _, v := range f.results {
        fmt.Println(v.ip + "," + v.start.Format(layout) + "," + v.mac + "," + v.hostname)
    }
}

func (f *leaseFinder)AllLeases() []Lease {
    ret := make([]Lease, 0, len(f.results))
    for _, v := range f.results {
        ret = append(ret, v)
    }
    return ret
}

func AllHostname(leases []Lease) string {
    var ret string
    for _, v := range leases {
        ret += v.hostname + "\n"
    }
    ret = strings.TrimRight(ret, "\n")
    return ret
}
