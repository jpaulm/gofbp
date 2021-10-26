package core

import (
	"fmt"
	"sync"
)

type networkStatus struct {
	mu              sync.Mutex
	deadlocked      bool
	running         int32
	liveConnections int32
	suspended       int32
}

func (network *Network) incLiveConnection() {
	network.status.mu.Lock()
	defer network.status.mu.Unlock()
	network.status.liveConnections++
}

func (network *Network) decLiveConnection() {
	network.status.mu.Lock()
	defer network.status.mu.Unlock()
	network.status.liveConnections--
}

func (network *Network) processTransitioned(from, to ProcessStatus) {
	if from == to {
		return
	}

	type tx [2]ProcessStatus

	switch (tx{from, to}) {
	case tx{NotStarted, Dormant}:
		network.status.mu.Lock()
		network.status.running++
		network.status.mu.Unlock()

	// these have no impact on deadlock detection
	case tx{Dormant, Active}:
	case tx{Active, Dormant}:

	case tx{Active, SuspendedSend}, tx{Active, SuspendedRecv}:
		network.status.mu.Lock()
		network.status.suspended++
		network.status.running--
		leftRunning := network.status.running
		network.status.mu.Unlock()

		if leftRunning <= 0 && network.status.liveConnections == 0 {
			network.deadlockDetected()
		}

	case tx{SuspendedSend, Active}, tx{SuspendedRecv, Active}:
		network.status.mu.Lock()
		network.status.suspended--
		network.status.running++
		network.status.mu.Unlock()

	case tx{Dormant, Terminated}, tx{Active, Terminated}:
		network.status.mu.Lock()
		network.status.running--
		running, suspended := network.status.running, network.status.suspended
		network.status.mu.Unlock()

		if running == 0 && suspended > 0 && network.status.liveConnections == 0 {
			network.deadlockDetected()
		}

	default:
		network.status.mu.Lock()
		defer network.status.mu.Unlock()

		if !network.status.deadlocked {
			panic(fmt.Sprintf("unhandled transition %s -> %s", from, to))
		}
	}
}

func (network *Network) deadlockDetected() {
	network.status.mu.Lock()
	defer network.status.mu.Unlock()

	if network.status.deadlocked {
		// avoid printing status twice
		return
	}
	network.status.deadlocked = true

	fmt.Println("\nDeadlock detected!")
	for _, proc := range network.procs {
		fmt.Printf("\t%-15s\t%s\n", proc.name, proc.status())
	}
	fmt.Println(network.goroutineTrace())
	panic("Deadlock!")
}
