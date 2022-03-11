package testing

import (
	"testing"

	"github.com/jpaulm/gofbp/components/testrtn"
	"github.com/jpaulm/gofbp/core"
	"github.com/jpaulm/gofbp/testing/components"
)

func TestLoadBal(t *testing.T) {
	params, err := core.LoadXMLParams("../params.xml")
	if err != nil {
		panic(err)
	}
	net := core.NewNetwork("TestLoadBal")
	net.SetParams(params)

	sender := net.NewProc("Sender", &testrtn.Sender{})

	proc1 := net.NewProc("DispFields", &testrtn.Marshal{})

	proc2 := net.NewProc("LoadBalance", &testrtn.LoadBalance{})

	proc3a := net.NewProc("Receiver0", &testrtn.Receiver{})
	proc3b := net.NewProc("Receiver1", &components.DelayedReceiver{})
	proc3c := net.NewProc("Receiver2", &components.DelayedReceiver{})

	net.Initialize("40", sender, "COUNT")
	net.Connect(sender, "OUT", proc1, "IN", 6)
	net.Connect(proc1, "OUT", proc2, "IN", 6)
	net.Connect(proc2, "OUT[0]", proc3a, "IN", 6)
	net.Connect(proc2, "OUT[1]", proc3b, "IN", 6)
	net.Connect(proc2, "OUT[2]", proc3c, "IN", 6)

	net.Run()
}
