package testing

import (
	"testing"

	"github.com/jpaulm/gofbp/components/io"
	"github.com/jpaulm/gofbp/components/testrtn"
	"github.com/jpaulm/gofbp/core"
)

func TestIntSender(t *testing.T) {
	net := core.NewNetwork("TestIntSender", nil)

	proc1 := net.NewProc("IntSender", &testrtn.IntSender{})
	proc2 := net.NewProc("WriteToFile", &io.WriteFile{})
	proc3 := net.NewProc("WriteToConsole", &testrtn.WriteToConsole{})

	net.Initialize("40", proc1, "COUNT")
	net.Initialize("numbers.txt", proc2, "FILENAME")
	net.Connect(proc1, "OUT", proc2, "IN", 6)
	net.Connect(proc2, "OUT", proc3, "IN", 6)
	net.Run()
}
