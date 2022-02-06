package testing

import (
	"testing"

	"github.com/jpaulm/gofbp/components/io"
	"github.com/jpaulm/gofbp/components/testrtn"
	"github.com/jpaulm/gofbp/core"
)

func TestInfQueueAsMain(t *testing.T) {
	params, err := core.LoadXMLParams("../params.xml")
	if err != nil {
		panic(err)
	}
	net := core.NewNetwork("TestInfQueueAsMain")
	net.SetParams(params)

	proc1 := net.NewProc("Sender", &testrtn.Sender{})
	proc2 := net.NewProc("WriteFile", &io.WriteFile{})
	proc3 := net.NewProc("ReadFile", &io.ReadFile{})
	proc4 := net.NewProc("WriteToConsole", &testrtn.WriteToConsole{})

	net.Initialize("40", proc1, "COUNT")
	net.Initialize("testdata/infqueue", proc2, "FILENAME")
	net.Initialize("testdata/infqueue", proc3, "FILENAME")
	net.Connect(proc1, "OUT", proc2, "IN", 6)
	net.Connect(proc2, "*", proc3, "*", 6)
	net.Connect(proc3, "OUT", proc4, "IN", 6)
	net.Run()
}
