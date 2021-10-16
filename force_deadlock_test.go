package main

import (
	"testing"

	"github.com/jpaulm/gofbp/components/testrtn"
	"github.com/jpaulm/gofbp/core"
)

func TestForceDeadlock(t *testing.T) {
	//t.Skip("Deadlock detection is designed to crash!")

	net := core.NewNetwork("ForceDeadlock")

	proc1 := net.NewProc("Sender", &testrtn.Sender{})
	proc2 := net.NewProc("Counter", &testrtn.Counter{})
	proc3 := net.NewProc("Concat", &testrtn.ConcatStr{})

	proc4 := net.NewProc("WriteToConsole", &testrtn.WriteToConsole{})

	net.Initialize("15", proc1, "COUNT")

	net.Connect(proc1, "OUT", proc2, "IN", 6)
	net.Connect(proc2, "OUT", proc3, "IN[1]", 6)
	net.Connect(proc2, "COUNT", proc3, "IN[0]", 6)
	net.Connect(proc3, "OUT", proc4, "IN", 6)

	net.Run()
}
