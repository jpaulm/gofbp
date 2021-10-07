package core

type Component interface {
	Setup(*Process)
	Execute(*Process)
}

type ComponentWithMustRun interface {
	Component
	MustRun() bool
}
