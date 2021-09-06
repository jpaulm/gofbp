package main

import (
	"github.com/jpaulm/gofbp/components/test-rtns/receiver"
	"github.com/jpaulm/gofbp/components/test-rtns/sender"
	"github.com/jpaulm/gofbp/core"
)

func main() {
	var net *core.Network = core.NewNetwork("test_net")

	proc1 := net.NewProc("Sender", sender.Execute, sender.OpenPorts)

	proc1a := net.NewProc("Sender2", sender.Execute, sender.OpenPorts)

	proc2 := net.NewProc("Receiver", receiver.Execute, receiver.OpenPorts)

	net.Connect(proc1, "OUT", proc2, "IN", 6)
	net.Connect(proc1a, "OUT", proc2, "IN", 6)

	net.Run()
}
