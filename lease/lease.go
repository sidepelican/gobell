package lease

import (
    "time"
    "strings"
)

type Lease struct {
    Start    *time.Time
    Ip       string
    Mac      string
    Hostname string
}

type Leases []Lease

func (leases Leases)AllHostname() string {
    var ret string
    for _, v := range leases {
        ret += v.Hostname + "\n"
    }
    ret = strings.TrimRight(ret, "\n")
    return ret
}
