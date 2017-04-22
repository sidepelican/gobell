package lease

import (
    "testing"
    "fmt"
)

func TestParse(t *testing.T) {

    leases, err := Parse("dhcpd_sample.leases")
    if err != nil {
        t.Fatal(err)
    }

    for _, l := range leases {
        fmt.Println(l)
    }
}
