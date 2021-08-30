package main

import (
	"runtime"

	core "github.com/jpaulm/gofbp/src/core"
)

var cc chan int = make(chan int, 10)

func main() {
	runtime.GOMAXPROCS(4)

	var net *core.Network = core.NewNetwork("test_net")

	net.Run()
}
