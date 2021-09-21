package main

import (
	"github.com/jpaulm/gofbp/components/io"
	"github.com/jpaulm/gofbp/components/testrtn"
	"github.com/jpaulm/gofbp/core"
)

// Merge application

func main() {

	// runtime.GOMAXPROCS(16)

	var net *core.Network = core.NewNetwork("test_net")

	proc1 := net.NewProc("ReadFile", &io.ReadFile{})

	proc2 := net.NewProc("WriteToConsole", &testrtn.WriteToConsole{})

	net.Initialize("C:\\Users\\Paul\\Documents\\GitHub\\gofbp\\.project", proc1, "FILENAME")
	net.Connect(proc1, "OUT", proc2, "IN", 6)

	net.Run()
}
