package components

import (
	"fmt"
	"reflect"

	"github.com/jpaulm/gofbp/core"
)

var Name string = "Receiver"

func Execute(p *core.Process) {
	fmt.Println("Receiver started")
	for {

		var pkt = p.Receive(p.InConn)
		if pkt == nil {
			break
		}
		v := reflect.ValueOf(pkt.Contents) // display contents - assume string
		s := v.String()
		fmt.Println("Output: " + s)
	}
	fmt.Println("Receiver ended")
}
