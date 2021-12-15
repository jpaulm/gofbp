package main

// Thiis differs from TestSubnet1 because it uses SSSubnet1 (SubstreamSenstive) instead of simple Subnet1...

import (
	"testing"

	"github.com/tyoung3/gofbp"
	"github.com/tyoung3/gofbp/subnets"
	"github.com/tyoung3/gofbp/testrtn"
)

func TestSubnet2(t *testing.T) {
	net := gofbp.NewNetwork("TestSubnet2")

	proc1 := net.NewProc("SubstreamSender", &testrtn.SubstreamSender{}) // sends multiple substreams of varying lengths

	proc2 := net.NewProc("RunSSSubnet", &subnets.SSSubnet1{})

	proc3 := net.NewProc("WriteToConsole2", &testrtn.WriteToConsole{})

	net.Initialize("20", proc1, "COUNT")
	net.Connect(proc1, "OUT", proc2, "IN", 6)
	net.Connect(proc2, "OUT", proc3, "IN", 6)

	net.Run()
}
