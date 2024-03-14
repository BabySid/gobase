package main

import (
	"fmt"
	"github.com/BabySid/gobase/log_sub"
	"log"
)

func main() {
	con, err := log_sub.NewConsumer(log_sub.Config{
		Location: &log_sub.SeekInfo{
			FileName: "20240314.log",
			Offset:   0,
			Whence:   0,
		},
		DateTimeLogLayout: &log_sub.DateTimeLayout{Layout: "20060102.log"},
	})

	if err != nil {
		panic(err)
	}

	i := 0
	for {
		line := <-con.Lines
		if line.Err != nil {
			log.Fatalf("got an error from consumeLog: err=%v", err)
			continue
		}
		fmt.Println(i, line)
		i += 1
	}
}
