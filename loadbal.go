package main

import (
	"github.com/jpaulm/gofbp/components/testrtn"

	"github.com/jpaulm/gofbp/core"
)

func main() {
	net := core.NewNetwork("loadbal", nil)

	receiver2 := net.NewProc("Receiver2", &testrtn.Receiver{})
	receiver1 := net.NewProc("Receiver1", &testrtn.Receiver{})
	receiver0 := net.NewProc("Receiver0", &testrtn.DelayedReceiver{})
	loadBalance := net.NewProc("LoadBalance", &testrtn.LoadBalance{})
	sender := net.NewProc("Sender", &testrtn.Sender{})

	net.Initialize("40", sender, "COUNT")
	net.Connect(loadBalance, "OUT[2]", receiver2, "IN", 6)
	net.Connect(loadBalance, "OUT[1]", receiver1, "IN", 6)
	net.Connect(loadBalance, "OUT[0]", receiver0, "IN", 6)
	net.Connect(sender, "OUT", loadBalance, "IN", 6)

	net.Run()
}
