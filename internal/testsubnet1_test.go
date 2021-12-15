package main

import (
	"testing"

	"github.com/jpaulm/gofbp"
	"github.com/jpaulm/gofbp/subnets"
	"github.com/jpaulm/gofbp/testrtn"
)

func TestSubnet1(t *testing.T) {
	net := gofbp.NewNetwork("TestSubnet1")

	proc1 := net.NewProc("Sender1", &testrtn.Sender{})

	proc2 := net.NewProc("RunSubnet", &subnets.Subnet1{})

	proc3 := net.NewProc("WriteToConsole2", &testrtn.WriteToConsole{})

	net.Initialize("15", proc1, "COUNT")
	net.Connect(proc1, "OUT", proc2, "IN", 6)
	net.Connect(proc2, "OUT", proc3, "IN", 6)

	net.Run()
}
