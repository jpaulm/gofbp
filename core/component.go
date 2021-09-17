package core

type Component interface {
	Setup(*Process)
	Execute(*Process)
	//GetMustRun(*Process) bool
}
