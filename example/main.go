package main

import (
	"fmt"
	"github.com/BabySid/gobase"
	"os"
	"syscall"
	"time"
)

func main() {
	ss := gobase.NewSignalSet()
	ss.Register(syscall.SIGTERM, exit)

	for {
		time.Sleep(time.Second * 10)
		fmt.Printf("%s is alive\n", os.Args[0])
	}
}

func exit(sig os.Signal) {
	fmt.Printf("%s exit by recving the signal %v\n", os.Args[0], sig)
	os.Exit(0)
}
