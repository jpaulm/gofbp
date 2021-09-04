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

	proc.OutConn = net.NewConnection(10)

	proc2 := net.NewProc(comp2.Execute)

	proc2.InConn = proc.OutConn

	net.Run()
}
