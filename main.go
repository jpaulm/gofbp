package main

import (
	"github.com/jpaulm/gofbp/components/testrtn"
	"github.com/jpaulm/gofbp/core"
)

func main() {
	var net *core.Network = core.NewNetwork("test_net")

	proc1 := net.NewProc("Sender", testrtn.NewSender())

	proc1a := net.NewProc("Sender2", testrtn.NewSender())

	proc2 := net.NewProc("Receiver", testrtn.NewReceiver())

	net.Connect(proc1, "OUT", proc2, "IN", 6)
	net.Connect(proc1a, "OUT", proc2, "IN", 6)

	net.Run()
}
