package main

import (
	"testing"

	"github.com/jpaulm/gofbp/components/testrtn"
	"github.com/jpaulm/gofbp/core"
)

func TestMerge(t *testing.T) {
	// runtime.GOMAXPROCS(16)

	var net *core.Network = core.NewNetwork("test_net")

	proc1 := net.NewProc("Sender1", &testrtn.Sender{})
	proc2 := net.NewProc("Sender2", &testrtn.Sender{})

	proc3 := net.NewProc("Receiver", &testrtn.Receiver{})

	net.Initialize("15", proc1, "COUNT")
	net.Initialize("10", proc2, "COUNT")
	net.Connect(proc1, "OUT", proc3, "IN", 6)
	net.Connect(proc2, "OUT", proc3, "IN", 6)

	net.Run()
}

func TestConcat(t *testing.T) {
	var net *core.Network = core.NewNetwork("test_net")

	proc1 := net.NewProc("Sender", &testrtn.Sender{})

	proc1a := net.NewProc("Sender2", &testrtn.Sender{})

	proc2 := net.NewProc("ConcatStr", &testrtn.ConcatStr{})

	proc3 := net.NewProc("Receiver", &testrtn.Receiver{})

	net.Initialize("15", proc1, "COUNT")
	net.Initialize("10", proc1a, "COUNT")
	net.Connect(proc1, "OUT", proc2, "IN[0]", 6)
	net.Connect(proc1a, "OUT", proc2, "IN[1]", 6)
	net.Connect(proc2, "OUT", proc3, "IN", 6)

	net.Run()
}

func TestRRDist(t *testing.T) {
	var net *core.Network = core.NewNetwork("test_net")

	proc1 := net.NewProc("Sender", &testrtn.Sender{})

	proc2 := net.NewProc("RoundRobinSender", &testrtn.RoundRobinSender{})

	proc3a := net.NewProc("Receiver1", &testrtn.Receiver{})
	proc3b := net.NewProc("Receiver2", &testrtn.Receiver{})
	proc3c := net.NewProc("Receiver3", &testrtn.Receiver{})

	net.Initialize("15", proc1, "COUNT")
	net.Connect(proc1, "OUT", proc2, "IN", 6)
	net.Connect(proc2, "OUT[0]", proc3a, "IN", 6)
	net.Connect(proc2, "OUT[1]", proc3b, "IN", 6)
	net.Connect(proc2, "OUT[2]", proc3c, "IN", 6)

	net.Run()
}
