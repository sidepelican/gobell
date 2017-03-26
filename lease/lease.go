package lease

import (
    "time"
    "strings"
)

type Lease struct {
    start    *time.Time
    ip       string
    mac      string
    hostname string
}

func AllHostname(leases []Lease) string {
    var ret string
    for _, v := range leases {
        ret += v.hostname + "\n"
    }
    ret = strings.TrimRight(ret, "\n")
    return ret
}
