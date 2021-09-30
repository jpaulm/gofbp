package main

import (
	"github.com/jpaulm/gofbp/components/io"
	//"github.com/jpaulm/gofbp/components/testrtn"
	"os"

	"github.com/jpaulm/gofbp/core"
)

// Concat

func main4() {

	net := core.NewNetwork("CopyFile")

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
