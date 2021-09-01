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

	//var p = comp.Process

	var sender core.Component = comp.Sender
	proc := net.NewProc("Sender", sender)

	//var myFun = proc.Execute(*core.Process)
	//proc.ProcFun = myFun
	proc.OutConn = net.NewConnection()
	net.Wg.Add(1)
	go proc.Run(net.Wg)

	net.Wg.Wait()

	net.Run()
}
