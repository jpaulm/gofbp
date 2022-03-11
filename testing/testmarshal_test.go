package testing

import (
	"testing"

	"github.com/jpaulm/gofbp/components/testrtn"
	"github.com/jpaulm/gofbp/core"
	"github.com/jpaulm/gofbp/testing/components"
	"github.com/jpaulm/gofbp/testing/structs"
)

func TestMarshal(t *testing.T) {
	params, err := core.LoadXMLParams("../params.xml")
	if err != nil {
		panic(err)
	}
	net := core.NewNetwork("TestMarshal")
	net.SetParams(params)

	name1 := &structs.Name{FirstName: "John", MidInit: "Q", LastName: "Smith"}
	emp1 := &structs.Employee{Name: name1, Age: 24, Salary: 344444}

	sender := net.NewProc("EmpSender", &components.EmpSender{})

	proc2 := net.NewProc("Marshal", &testrtn.Marshal{})

	proc3 := net.NewProc("Show", &testrtn.WriteToConsole{})

	net.Initialize("20", sender, "COUNT")
	net.Initialize(emp1, sender, "DATA")
	net.Connect(sender, "OUT", proc2, "IN", 6)
	net.Connect(proc2, "OUT", proc3, "IN", 6)
	net.Run()
}
