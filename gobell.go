package main

import (
    "bufio"
    "fmt"
    "io"
    "os"
    "regexp"
    "time"
    //"github.com/go-fsnotify/fsnotify"
)

type Lease struct {
    start    *time.Time
    ip       string
    mac      string
    hostname string
}

type LeaseFinder struct {
    results []Lease
    current *Lease
}

func (f *LeaseFinder)findStart(ip string) {

    if f.current != nil {
        fmt.Println("findStartLease called before call findEnd. something wrong")
        f.current = nil
    }

    f.current = &Lease{ip: ip}
}

func (f *LeaseFinder)findStartTime(startString string) {

    const layout = "2006/01/02 15:04:05"
    t, err := time.Parse(layout, startString)
    if err != nil {
        fmt.Println(err)
        return
    }

    f.current.start = &t
}

func (f *LeaseFinder)findMac(mac string) {
    f.current.mac = mac
}

func (f *LeaseFinder)findHostname(hostname string) {
    f.current.hostname = hostname
}

func (f *LeaseFinder)findEnd() {

    defer func() {
        f.current = nil
    }()

    // something dropped
    if f.current.hostname == "" || f.current.start == nil || f.current.ip == "" || f.current.mac == "" {
        return
    }

    // all values are completed
    f.results = append(f.results, *f.current)

    return
}

func (f *LeaseFinder)printAll() {

    const layout = "2006/01/02 15:04:05"
    for _, v := range f.results {
        fmt.Println(v.ip + "," + v.start.Format(layout) + "," + v.mac + "," + v.hostname)
    }
}

func main() {

    fp, err := os.Open("sample/dhcpd.leases")
    if err != nil {
        panic(err)
    }
    defer fp.Close()

    // regexp setup
    rLease := regexp.MustCompile(`lease ([0-9¥.]+) {`)
    rStarts := regexp.MustCompile(`starts ([0-9]) (.*);`)
    rHwEth := regexp.MustCompile(`hardware ethernet ([0-9A-Fa-f:-]+);`)
    rHostname := regexp.MustCompile(`client-hostname "(.*)";`)
    rEnd := regexp.MustCompile(`}`)

    lf := LeaseFinder{}

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
            lf.findStart(res[len(res)-1])
        }
        res = rStarts.FindStringSubmatch(line)
        if res != nil {
            lf.findStartTime(res[len(res)-1])
        }
        res = rHwEth.FindStringSubmatch(line)
        if res != nil {
            lf.findMac(res[len(res)-1])
        }
        res = rHostname.FindStringSubmatch(line)
        if res != nil {
            lf.findHostname(res[len(res)-1])
        }
        if rEnd.MatchString(line) {
            lf.findEnd()
        }
    }

    lf.printAll()
}