package main

import (
	"github.com/jpaulm/gofbp/components/testrtn"
	"github.com/jpaulm/gofbp/core"
)

func main2() {
	var net *core.Network = core.NewNetwork("test_net")

	proc1 := net.NewProc("Sender", &testrtn.Sender{})

	proc1a := net.NewProc("Sender2", &testrtn.Sender{})

	proc2 := net.NewProc("Receiver", &testrtn.Receiver{})

	net.Initialize("15", proc1, "COUNT")
	net.Initialize("10", proc1a, "COUNT")
	net.Connect(proc1, "OUT", proc2, "IN", 6)
	net.Connect(proc1a, "OUT", proc2, "IN", 6)

	net.Run()
}
