/*Package core implements gofbp's run time engine.*/
package core

/*The Component interface implements component Setup and Execute methods.*/
type Component interface {
	Setup(*Process)
	Execute(*Process)
}

/*The ComponentWithMustRun interface implements the MustRun method.*/
type ComponentWithMustRun interface {
	Component
	MustRun() bool
}
