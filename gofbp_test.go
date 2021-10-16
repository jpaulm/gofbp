package main

import (
	"testing"

	"path/filepath"

	"github.com/jpaulm/gofbp/components/io"
	"github.com/jpaulm/gofbp/components/testrtn"
	"github.com/jpaulm/gofbp/core"
)

func TestMerge(t *testing.T) {
	net := core.NewNetwork("Merge")

	proc1 := net.NewProc("Sender1", &testrtn.Sender{})
	proc2 := net.NewProc("Sender2", &testrtn.Sender{})

	proc3 := net.NewProc("WriteToConsole", &testrtn.WriteToConsole{})

	net.Initialize("15", proc1, "COUNT")
	net.Initialize("10", proc2, "COUNT")
	net.Connect(proc1, "OUT", proc3, "IN", 6)
	net.Connect(proc2, "OUT", proc3, "IN", 6)

	net.Run()
}

func TestConcat(t *testing.T) {
	net := core.NewNetwork("Concat")

	proc1 := net.NewProc("Sender", &testrtn.Sender{})

	proc1a := net.NewProc("Sender2", &testrtn.Sender{})

	proc2 := net.NewProc("ConcatStr", &testrtn.ConcatStr{})

	proc3 := net.NewProc("WriteToConsole", &testrtn.WriteToConsole{})

	net.Initialize("15", proc1, "COUNT")
	net.Initialize("10", proc1a, "COUNT")
	net.Connect(proc1, "OUT", proc2, "IN[0]", 6)
	net.Connect(proc1a, "OUT", proc2, "IN[1]", 6)
	net.Connect(proc2, "OUT", proc3, "IN", 6)

	net.Run()
}

func TestRRDist(t *testing.T) {
	net := core.NewNetwork("RRDist")

	proc1 := net.NewProc("Sender", &testrtn.Sender{})

	proc2 := net.NewProc("RoundRobinSender", &testrtn.RoundRobinSender{})

	proc3a := net.NewProc("WriteToConsole", &testrtn.WriteToConsole{})
	proc3b := net.NewProc("Receiver1", &testrtn.Receiver{})
	proc3c := net.NewProc("Receiver2", &testrtn.Receiver{})

	net.Initialize("15", proc1, "COUNT")
	net.Connect(proc1, "OUT", proc2, "IN", 6)
	net.Connect(proc2, "OUT[0]", proc3a, "IN", 6)
	net.Connect(proc2, "OUT[1]", proc3b, "IN", 6)
	net.Connect(proc2, "OUT[2]", proc3c, "IN", 6)

	net.Run()
}

func TestCopyFile(t *testing.T) {
	net := core.NewNetwork("CopyFile")

	proc1 := net.NewProc("ReadFile", &io.ReadFile{})

	proc2 := net.NewProc("WriteFile", &io.WriteFile{})

	net.Initialize(filepath.Join("testdata", "testdata.txt"), proc1, "FILENAME")
	net.Initialize(filepath.Join("testdata", "copy-file.tmp"), proc2, "FILENAME")
	net.Connect(proc1, "OUT", proc2, "IN", 6)

	net.Run()
}

func TestDoSelect1(t *testing.T) {

	// port REJ from proc2 is connected

	net := core.NewNetwork("DoSelect")

	proc1 := net.NewProc("ReadFile", &io.ReadFile{})
	proc2 := net.NewProc("Select", &testrtn.Selector{})
	proc3a := net.NewProc("WriteFile", &io.WriteFile{})
	proc3b := net.NewProc("WriteToConsole", &testrtn.WriteToConsole{})

	net.Initialize(filepath.Join("testdata", "testdata.txt"), proc1, "FILENAME")
	net.Initialize("X", proc2, "PARAM")
	net.Initialize(filepath.Join("testdata", "do-select.tmp"), proc3a, "FILENAME")
	net.Connect(proc1, "OUT", proc2, "IN", 6)
	net.Connect(proc2, "ACC", proc3a, "IN", 6)
	net.Connect(proc2, "REJ", proc3b, "IN", 6)

	net.Run()
}

func TestDoSelect2(t *testing.T) {

	// port REJ from proc2 is NOT connected

	net := core.NewNetwork("DoSelect2")

	proc1 := net.NewProc("ReadFile", &io.ReadFile{})
	proc2 := net.NewProc("Select", &testrtn.Selector{})
	proc3a := net.NewProc("WriteFile", &io.WriteFile{})
	//proc3b := net.NewProc("WriteToConsole", &testrtn.WriteToConsole{})

	net.Initialize(filepath.Join("testdata", "testdata.txt"), proc1, "FILENAME")
	net.Initialize("X", proc2, "PARAM")
	net.Initialize(filepath.Join("testdata", "do-select.tmp"), proc3a, "FILENAME")
	net.Connect(proc1, "OUT", proc2, "IN", 6)
	net.Connect(proc2, "ACC", proc3a, "IN", 6)
	//net.Connect(proc2, "REJ", proc3b, "IN", 6)

	net.Run()
}
func TestWriteToConsUsingNL(t *testing.T) {

	net := core.NewNetwork("MergeToCons")

	proc1 := net.NewProc("Sender1", &testrtn.Sender{})
	proc2 := net.NewProc("Sender2", &testrtn.Sender{})

	proc3 := net.NewProc("WriteToConsNL", &testrtn.WriteToConsNL{})

	net.Initialize("15", proc1, "COUNT")
	net.Initialize("15", proc2, "COUNT")
	net.Connect(proc1, "OUT", proc3, "IN", 6)
	net.Connect(proc2, "OUT", proc3, "IN", 6)

	net.Run()
}
