package core

type Component interface {
	OpenPorts(*Process)
	Execute(*Process)
	//GetMustRun(*Process) bool
}
