// package ans (ws_ans_req.go)  is a totally .... intended just to fill out the structure of a GoFBP web socket server...
package ans

import (
	"testing"

	"github.com/jpaulm/gofbp/components/websocket"
	"github.com/jpaulm/gofbp/core"
)

func TestWebSocket(t *testing.T) {
	params, err := core.LoadXMLParams("params.xml")
	if err != nil {
		panic(err)
	}
	path := "localhost:8080"
	net := core.NewNetwork("TestWebSocket")
	net.SetParams(params)
	//proc1 := net.NewProc("WSRequest", &websocket.WSRequest{})
	proc1 := net.NewProc("WSRequest", &websocket.WSRequest{}) // Assumed to be working
	proc2 := net.NewProc("AnsReq", &WSAnsReq{})               // Testing this one.
	proc3 := net.NewProc("WSRespond", &websocket.WSRespond{}) // Assumed to be working
	net.Initialize(path, proc1, "ADDR")
	net.Connect(proc1, "OUT", proc2, "IN", 6)
	net.Connect(proc2, "OUT", proc3, "IN", 6)
	net.Run()
}
