package main

import (
	//"os"

	//"github.com/jpaulm/gofbp/components/io"
	"github.com/jpaulm/gofbp/components/testrtn"
	"github.com/jpaulm/gofbp/core"
	// "runtime"
)

func main() {

	net := core.NewNetwork("RRDist")

	proc1 := net.NewProc("Sender", &testrtn.Sender{})

	proc2 := net.NewProc("RoundRobinSender", &testrtn.RoundRobinSender{})

	proc3a := net.NewProc("WriteToConsole", &testrtn.WriteToConsole{})
	proc3b := net.NewProc("Receiver1", &testrtn.Receiver{})
	proc3c := net.NewProc("Receiver2", &testrtn.Receiver{})

	net.Initialize("15", proc1, "COUNT")
	net.Connect(proc1, "OUT", proc2, "IN", 6)
	net.Connect(proc2, "OUT[0]", proc3a, "IN", 6)
	net.Connect(proc2, "OUT[1]", proc3b, "IN", 6)
	net.Connect(proc2, "OUT[2]", proc3c, "IN", 6)

	net.Run()

}
