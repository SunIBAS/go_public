package main

import (
	"fmt"
	"github.com/go-vgo/robotgo"
	"time"
)

var ch = make(chan string)

func main() {
	go listenEnter()
	go listenQ()
	for {
		time.Sleep(time.Millisecond * 100)
		c := <-ch
		fmt.Println(c)
		if c == "q" {
			go listenQ()
		} else if c == "s" {
			go listenEnter()
		}
	}
	fmt.Scanln()
}
func listenQ() {
	ok := robotgo.AddEvents("q")
	if ok {
		//fmt.Println("q")
		ch <- "q"
	}
}
func listenEnter() {
	ok := robotgo.AddEvents("s")
	if ok {
		//fmt.Println("s")
		ch <- "s"
	}
}
