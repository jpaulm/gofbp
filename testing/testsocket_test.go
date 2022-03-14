package testing

import (
	"testing"

	"github.com/jpaulm/gofbp/components/io"
	"github.com/jpaulm/gofbp/components/testrtn"
	"github.com/jpaulm/gofbp/core"
)

func TestSocket(t *testing.T) {
	params, err := core.LoadXMLParams("../params.xml")
	if err != nil {
		panic(err)
	}
	net := core.NewNetwork()
	net.SetParams(params)
	proc1 := net.NewProc("Sender", &testrtn.Sender{})
	proc2 := net.NewProc("WriteSocket", &io.WriteToSocket{})
	proc3 := net.NewProc("ReadSocket", &io.ReadFromSocket{})
	proc4 := net.NewProc("WriteToConsole", &testrtn.WriteToConsole{})

	net.Initialize("20", proc1, "COUNT")
	net.Initialize("192.168.0.0:4444", proc2, "PORT")
	net.Initialize("192.168.0.0:4444", proc3, "PORT")
	net.Connect(proc1, "OUT", proc2, "IN", 6)
	net.Connect(proc3, "OUT", proc4, "IN", 6)
	net.Run()
}
