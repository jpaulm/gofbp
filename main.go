package main

import (
	comp "github.com/jpaulm/gofbp/components/sender"
	"github.com/jpaulm/gofbp/core"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(4)

	var net *core.Network = core.NewNetwork("test_net")

	proc := net.NewProc(comp.Execute)

	proc.OutConn = net.NewConnection()

	net.Run()
}
