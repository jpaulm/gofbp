package testing

import (
	"testing"

	"github.com/jpaulm/gofbp/components/testrtn"
	"github.com/jpaulm/gofbp/core"
)

func TestDropOldest(t *testing.T) {
	params, err := core.LoadXMLParams("../params.xml")
	if err != nil {
		panic(err)
	}
	net := core.NewNetwork("TestDropOldest")
	net.SetParams(params)
	proc1 := net.NewProc("IntSender", &testrtn.IntSender{})
	proc2 := net.NewProc("DelayedReceiver", &testrtn.DelayedReceiver{})

	net.Initialize("50", proc1, "COUNT")
	conn1 := net.Connect(proc1, "OUT", proc2, "IN", 6)
	net.SetDropOldest(conn1)
	net.Run()
}
