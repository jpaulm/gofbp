package testing

import (
	"testing"

	"github.com/jpaulm/gofbp/components/subnets"
	"github.com/jpaulm/gofbp/components/testrtn"
	"github.com/jpaulm/gofbp/core"
)

func TestSubnet1(t *testing.T) {
	params, err := core.LoadXMLParams("../params.xml")
	if err != nil {
		panic(err)
	}
	net := core.NewNetwork("TestSubnet1")
	net.SetParams(params)

	proc1 := net.NewProc("Sender1", &testrtn.Sender{})

	proc2 := net.NewProc("RunSubnet", &subnets.Subnet1{})

	proc3 := net.NewProc("WriteToConsole", &testrtn.WriteToConsole{})

	net.Initialize("15", proc1, "COUNT")
	net.Connect(proc1, "OUT", proc2, "IN", 6)
	net.Connect(proc2, "OUT", proc3, "IN", 6)

	net.Run()
}
