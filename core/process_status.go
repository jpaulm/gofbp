package core

import "fmt"

type ProcessStatus int32

const (
	NotStarted ProcessStatus = iota
	Active
	Dormant
	SuspendedSend
	SuspendedRecv
	Terminated
)

func (status ProcessStatus) String() string {
	switch status {
	case NotStarted:
		return "NotStarted"
	case Active:
		return "Active"
	case Dormant:
		return "Dormant"
	case SuspendedSend:
		return "SuspendedSend"
	case SuspendedRecv:
		return "SuspendedRecv"
	case Terminated:
		return "Terminated"
	default:
		return fmt.Sprintf("ProcessStatus(%d)", status)
	}
}
