package gofbp

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
	SetArrayItem(c *OutPort, i int)
	ArrayLength() int
	GetItemWithFewestIPs() int
}
