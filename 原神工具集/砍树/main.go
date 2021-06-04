package main

import (
	"C"
	"fmt"
	"github.com/go-vgo/robotgo"
)

func main() {
	//robotgo.KeyTap(os.Args[1])
	//robotgo.MoveMouse(0,0)
	//robotgo.Sleep(2000)
	x, y := robotgo.GetMousePos()
	fmt.Println(x)
	fmt.Println(y)
	//robotgo.MoveMouse(200,y)
	robotgo.MoveMouseSmooth(x-100, y)
}
