package main

import (
	"github.com/jpaulm/gofbp/components/testrtn"
	"github.com/jpaulm/gofbp/core"
	// "runtime"
)

func main() {

	// runtime.GOMAXPROCS(16)

	var net *core.Network = core.NewNetwork("test_net")

	proc1 := net.NewProc("Sender", &testrtn.Sender{})

	proc2 := net.NewProc("RoundRobinSender", &testrtn.RoundRobinSender{})

	proc3a := net.NewProc("Receiver", &testrtn.Receiver{})
	proc3b := net.NewProc("Receiver", &testrtn.Receiver{})
	proc3c := net.NewProc("Receiver", &testrtn.Receiver{})

	net.Initialize("15", proc1, "COUNT")
	net.Connect(proc1, "OUT", proc2, "IN", 6)
	net.Connect(proc2, "OUT[0]", proc3a, "IN", 6)
	net.Connect(proc2, "OUT[1]", proc3b, "IN", 6)
	net.Connect(proc2, "OUT[2]", proc3c, "IN", 6)

	net.Run()
}
