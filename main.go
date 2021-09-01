package main

import (
	"runtime"

	comp "github.com/jpaulm/gofbp/components/sender"
	"github.com/jpaulm/gofbp/core"
)

//var cc chan int = make(chan int, 10)

func main() {
	runtime.GOMAXPROCS(4)

	var net *core.Network = core.NewNetwork("test_net")

	var sender core.Component = comp.Sender
	sender.Name = "Sender"

	proc := net.NewProc(sender)

	proc.OutConn = net.NewConnection()
	net.Wg.Add(1)
	go proc.Run(net.Wg)

	net.Wg.Wait()

	net.Run()
}
