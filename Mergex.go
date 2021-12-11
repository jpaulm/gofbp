package testing //change package name, or delete statement, if desired
// Generated (mostly) - however, had to change name to Mergex, as there is already a Merge in gofbp_test.go
import (
	"testing"
	//...
	"github.com/jpaulm/gofbp/components/testrtn" // not generated (yet)
	"github.com/jpaulm/gofbp/core"
)

func TestMergex(t *testing.T) {
	net := core.NewNetwork("Merge")
	sender2 := net.NewProc("Sender2", &testrtn.Sender{})
	write___to_console := net.NewProc("Write___To_Console", &testrtn.WriteToConsole{})
	sender1 := net.NewProc("Sender1", &testrtn.Sender{})
	net.Connect(sender1, "OUT", write___to_console, "IN", 6)
	net.Initialize("15", sender1, "IN")
	net.Initialize("10", sender2, "IN")
	net.Connect(sender2, "OUT", write___to_console, "IN", 6)
	net.Run()
}
