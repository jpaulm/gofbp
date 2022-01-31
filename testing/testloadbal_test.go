package testing

import (
	"testing"

	"github.com/jpaulm/gofbp/components/testrtn"
	"github.com/jpaulm/gofbp/core"
)

func TestLoadBal(t *testing.T) {
	net := core.NewNetwork("TestLoadBal")

	sender := net.NewProc("Sender", &testrtn.Sender{})

	proc2 := net.NewProc("LoadBalance", &testrtn.LoadBalance{})

	proc3a := net.NewProc("Receiver0", &testrtn.Receiver{})
	proc3b := net.NewProc("Receiver1", &testrtn.DelayedReceiver{})
	proc3c := net.NewProc("Receiver2", &testrtn.DelayedReceiver{})

	net.Initialize("40", sender, "COUNT")
	net.Connect(sender, "OUT", proc2, "IN", 6)
	net.Connect(proc2, "OUT[0]", proc3a, "IN", 6)
	net.Connect(proc2, "OUT[1]", proc3b, "IN", 6)
	net.Connect(proc2, "OUT[2]", proc3c, "IN", 6)

	net.Run()
}
