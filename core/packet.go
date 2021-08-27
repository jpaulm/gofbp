package main

type Packet struct {
	content interface{}
	pktType int
	owner   *Process
}
