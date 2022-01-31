package core

type outputCommon interface {
	Close()
	IsConnected() bool
}

type OutputConn interface {
	outputCommon
	send(*Process, *Packet) bool
	// IsConnected() bool
}

type OutputArrayConn interface {
	outputCommon
	GetArrayItem(i int) *OutPort
	setArrayItem(c *OutPort, i int)
	ArrayLength() int
	GetItemWithFewestIPs() int
}
