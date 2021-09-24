package main

import (
	"github.com/jpaulm/gofbp/components/io"
	//"github.com/jpaulm/gofbp/components/testrtn"
	"github.com/jpaulm/gofbp/core"
)

// Concat

func main() {

	var net *core.Network = core.NewNetwork("CopyFile")

	proc1 := net.NewProc("ReadFile", &io.ReadFile{})

	proc2 := net.NewProc("WriteFile", &io.WriteFile{})

	net.Initialize("C:\\Users\\Paul\\Documents\\GitHub\\gofbp\\.project", proc1, "FILENAME")
	net.Initialize("C:\\Users\\Paul\\Documents\\GitHub\\gofbp\\.project.copy", proc2, "FILENAME")
	net.Connect(proc1, "OUT", proc2, "IN", 6)

	net.Run()
}
