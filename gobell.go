package main

import (
	"bufio"
	"fmt"
	"os"
	"io"
	"regexp"

	//"github.com/go-fsnotify/fsnotify"
)

func main() {

	fp, err := os.Open("sample/dhcpd.leases")
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	// regexp setup
	rLease    := regexp.MustCompile(`lease [0-9¥.]+ {`)
    rStarts   := regexp.MustCompile(`starts .*;`)
    rHwEth    := regexp.MustCompile(`hardware ethernet [0-9A-Fa-f:-]+;`)
    rHostname := regexp.MustCompile(`client-hostname ".*";`)
    rEnd      := regexp.MustCompile(`}`)

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
        if rLease.MatchString(line) {
            fmt.Println("lease {")
        }
        if rStarts.MatchString(line) {
            fmt.Println("start: ")
        }
        if rHwEth.MatchString(line) {
            fmt.Println("hardware ethernet: ")
        }
        if rHostname.MatchString(line) {
            fmt.Println("client-hostname: ")
        }
        if rEnd.MatchString(line) {
            fmt.Println("}")
        }
	}
}
