package main

import (
	"bufio"
	"fmt"
	"os"
	"io"

	//"github.com/go-fsnotify/fsnotify"

	"io"
)

func main() {

	fmt.Println("Hello world!")

	fp, err := os.Open("sample/dhcpd.leases")
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	reader := bufio.NewReaderSize(fp, 4096)
	for line := ""; err == nil; line, err = reader.ReadString('\n') {
		fmt.Print(line)
	}
	if err != io.EOF {
		panic(err)
	}
}
