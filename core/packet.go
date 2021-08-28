package main

type Packet struct {
	contents interface{}
	pktType  int
	owner    *Process
}
