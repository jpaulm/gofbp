package main

import (
	"os"

	"github.com/jpaulm/gofbp/components/io"
	"github.com/jpaulm/gofbp/components/testrtn"
	"github.com/jpaulm/gofbp/core"
	// "runtime"
)

func main() {

	net := core.NewNetwork("DoSelect")

	proc1 := net.NewProc("ReadFile", &io.ReadFile{})
	proc2 := net.NewProc("Select", &testrtn.Selector{})
	proc3a := net.NewProc("WriteFile", &io.WriteFile{})
	proc3b := net.NewProc("WriteToConsole", &testrtn.WriteToConsole{})

	path, err := os.Getwd()
	if err != nil {
		panic("Can't find workspace directory")
	}

	net.Initialize(path+"\\testdata.txt", proc1, "FILENAME")
	net.Initialize("X", proc2, "PARAM")
	net.Initialize(path+"\\testdata.copy", proc3a, "FILENAME")
	net.Connect(proc1, "OUT", proc2, "IN", 6)
	net.Connect(proc2, "ACC", proc3a, "IN", 6)
	net.Connect(proc2, "REJ", proc3b, "IN", 6)

	net.Run()

}
