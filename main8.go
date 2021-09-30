package main

import (
	"github.com/jpaulm/gofbp/components/io"
	"github.com/jpaulm/gofbp/components/testrtn"
	"github.com/jpaulm/gofbp/core"
)

// Concat

func main8() {

	net := core.NewNetwork("DoSelect")

	proc1 := net.NewProc("ReadFile", &io.ReadFile{})
	proc2 := net.NewProc("Select", &testrtn.Selector{})
	proc3a := net.NewProc("WriteFile", &io.WriteFile{})
	proc3b := net.NewProc("WriteToConsole", &testrtn.WriteToConsole{})

	net.Initialize("C:\\Users\\Paul\\Documents\\GitHub\\gofbp\\.project", proc1, "FILENAME")
	net.Initialize("X", proc2, "PARAM")
	net.Initialize("C:\\Users\\Paul\\Documents\\GitHub\\gofbp\\.project.copy", proc3a, "FILENAME")
	net.Connect(proc1, "OUT", proc2, "IN", 6)
	net.Connect(proc2, "ACC", proc3a, "IN", 6)
	net.Connect(proc2, "REJ", proc3b, "IN", 6)

	net.Run()
}
