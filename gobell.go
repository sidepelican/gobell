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
	rLease    := regexp.MustCompile(`lease ([0-9¥.]+) {`)
    rStarts   := regexp.MustCompile(`starts ([0-9]) (.*);`)
    rHwEth    := regexp.MustCompile(`hardware ethernet ([0-9A-Fa-f:-]+);`)
    rHostname := regexp.MustCompile(`client-hostname "(.*)";`)
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
        var res []string = nil
        res = rLease.FindStringSubmatch(line)
        if res != nil {
            fmt.Println("lease " + res[len(res)-1] + " {")
        }
        res = rStarts.FindStringSubmatch(line)
        if res != nil {
            fmt.Println("start: " + res[len(res)-1])
        }
        res = rHwEth.FindStringSubmatch(line)
        if res != nil {
            fmt.Println("hardware ethernet: " + res[len(res)-1])
        }
        res = rHostname.FindStringSubmatch(line)
        if res != nil {
            fmt.Println("client-hostname: " + res[len(res)-1])
        }
        if rEnd.MatchString(line) {
            fmt.Println("}")
        }
	}
}
