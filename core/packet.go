package core

type Packet struct {
	Contents interface{}
	pktType  int
	owner    *Process
}
