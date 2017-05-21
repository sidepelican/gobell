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

func (u Leases) Len() int {
    return len(u)
}

func (u Leases) Swap(i, j int) {
    u[i], u[j] = u[j], u[i]
}

func (u Leases) Less(i, j int) bool {
    return u[i].Ip < u[j].Ip
}