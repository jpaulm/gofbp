package main

import (
	"runtime"
)

var cc chan int = make(chan int, 10)

func main() {
	runtime.GOMAXPROCS(4)

	var net *Network = NewNetwork("test_net")

	net.Run()
}
