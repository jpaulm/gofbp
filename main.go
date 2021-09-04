package main

import (
	"runtime"

	comp2 "github.com/jpaulm/gofbp/components/receiver"
	comp "github.com/jpaulm/gofbp/components/sender"
	"github.com/jpaulm/gofbp/core"
)

func main() {
	runtime.GOMAXPROCS(4)

	var net *core.Network = core.NewNetwork("test_net")

	proc := net.NewProc("Sender", comp.Execute)
	//proc.Name = "Sender"

	proc1a := net.NewProc("Sender2", comp.Execute)
	//proc1a.Name = "Sender2"

	proc.OutConn = net.NewConnection(6)

	proc2 := net.NewProc("Receiver", comp2.Execute) // Note different import key!
	//proc2.Name = "Receiver"

	proc2.InConn = proc.OutConn
	proc1a.OutConn = proc.OutConn // 2 outputs feeding 1 input (legal in FBP)

	proc.OutConn.UpStrmCnt = 2

	net.Run()
}
