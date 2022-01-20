package testing

import (
	"testing"

	"github.com/jpaulm/gofbp/components/websocket"

	"github.com/jpaulm/gofbp/components/testrtn"
	"github.com/jpaulm/gofbp/core"
)

func TestWebSocket(t *testing.T) {
	params, err := core.LoadXMLParams("../params.xml")
	if err != nil {
		panic(err)
	}
	net := core.NewNetwork("TestWebSocket")
	net.SetParams(params)
	proc1 := net.NewProc("WSServer", &websocket.WSServer{})
	proc2 := net.NewProc("WriteToConsole", &testrtn.WriteToConsole{})

	net.Connect(proc1, "OUT", proc2, "IN", 6)
	net.Run()
}
