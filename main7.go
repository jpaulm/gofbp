package main

import (
	"github.com/jpaulm/gofbp/components/testrtn"
	"github.com/jpaulm/gofbp/core"
)

// Concat

func main7() {

	var net *core.Network = core.NewNetwork("Concat")

	proc1 := net.NewProc("Sender", &testrtn.Sender{})

	proc1a := net.NewProc("Sender2", &testrtn.Sender{})

	proc2 := net.NewProc("ConcatStr", &testrtn.ConcatStr{})

	proc3 := net.NewProc("WriteToConsole", &testrtn.WriteToConsole{})

	net.Initialize("15", proc1, "COUNT")
	net.Initialize("10", proc1a, "COUNT")
	net.Connect(proc1, "OUT", proc2, "IN[0]", 6)
	net.Connect(proc1a, "OUT", proc2, "IN[1]", 6)
	net.Connect(proc2, "OUT", proc3, "IN", 6)

	net.Run()
}
