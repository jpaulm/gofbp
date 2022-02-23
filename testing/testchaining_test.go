package testing

import (
	"testing"

	//"github.com/jpaulm/gofbp/components/io"
	"github.com/jpaulm/gofbp/components/testrtn"
	"github.com/jpaulm/gofbp/core"
	"github.com/jpaulm/gofbp/testing/components"
)

func TestChaining(t *testing.T) {
	params, err := core.LoadXMLParams("../params.xml")
	if err != nil {
		panic(err)
	}
	net := core.NewNetwork()
	net.SetParams(params)
	proc1 := net.NewProc("ChainBuild", &components.ChainBuild{})
	proc2 := net.NewProc("Discard", &testrtn.Discard{})

	net.Connect(proc1, "OUT", proc2, "IN", 6)
	net.Run()
}
