package lease

import (
    "time"
    "strings"
)

type Lease struct {
    Start    *time.Time `json:"start"`
    Ip       string     `json:"ip"`
    Mac      string     `json:"mac"`
    Hostname string     `json:"hostname"`
}

type Leases []Lease

func (leases Leases)AllHostname() (ret string) {
    for _, v := range leases {
        ret += v.Hostname + "\n"
    }
    ret = strings.TrimRight(ret, "\n")
    return
}

func TrimMacAddr(s string) string {
    ret := strings.ToLower(s)
    ret = strings.Replace(ret,`-`, `:`, -1)
    return ret
}