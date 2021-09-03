package main

import (
	"reflect"
	"runtime"

	comp "github.com/jpaulm/gofbp/components/sender"
	"github.com/jpaulm/gofbp/core"
)

func getTypeName(t reflect.Type) string {
	return t.Name()
}

func main() {
	runtime.GOMAXPROCS(4)

	var net *core.Network = core.NewNetwork("test_net")

	proc := net.NewProc(comp.Execute)

	proc.OutConn = net.NewConnection()

	//net.Wg.Add(1)
	// go proc.Run(net)
	//net.Wg.Wait()

	net.Run()
}
