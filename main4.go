package main

import (
	"github.com/jpaulm/gofbp/components/io"
	//"github.com/jpaulm/gofbp/components/testrtn"
	"github.com/jpaulm/gofbp/core"
	"os"
)

// Concat

func main4() {

	var net *core.Network = core.NewNetwork("CopyFile")

	proc1 := net.NewProc("ReadFile", &io.ReadFile{})

	proc2 := net.NewProc("WriteFile", &io.WriteFile{})

	path, err := os.Getwd()
	if err != nil {
		panic("Can't find workspace directory")
	}

	net.Initialize(path+"\\testdata.txt", proc1, "FILENAME")
	net.Initialize(path+"\\testdata.copy", proc2, "FILENAME")
	net.Connect(proc1, "OUT", proc2, "IN", 6)

	net.Run()
}
