package main

import (
	"github.com/BabySid/gobase/log_sub"
)

func main() {
	_, err := log_sub.NewConsumer(log_sub.Config{
		Location: &log_sub.SeekInfo{
			FileName: "../log/2023060612.log",
			Offset:   0,
			Whence:   0,
		},
		DateTimeLogLayout: &log_sub.DateTimeLayout{Layout: "../log/2006010215.log"},
	})

	if err != nil {
		panic(err)
	}
}
