package main

import (
	comp2 "github.com/jpaulm/gofbp/components/receiver"
	comp "github.com/jpaulm/gofbp/components/sender"
	"github.com/jpaulm/gofbp/core"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(4)

	var net *core.Network = core.NewNetwork("test_net")

	proc := net.NewProc(comp.Execute)
	proc.Name = "Sender"

	proc1a := net.NewProc(comp.Execute)
	proc1a.Name = "Sender2"

	proc.OutConn = net.NewConnection(6)

	proc2 := net.NewProc(comp2.Execute)  
	proc2.Name = "Receiver"

	proc2.InConn = proc.OutConn
	proc1a.OutConn = proc.OutConn // 2 outputs feeding 1 input (legal in FBP)

	net.Run()
}
