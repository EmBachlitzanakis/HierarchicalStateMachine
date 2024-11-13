package main

import (
	"fmt"
	"time"
)

// Mocking HSM package structures and functions (placeholders)
type Event int

type Hsm struct {
	currentState State
}

type State interface {
	Init(hsm *Hsm)
	Event(hsm *Hsm, e Event) bool
}

func NewHsm(initialState State) *Hsm {
	hsm := &Hsm{currentState: initialState}
	initialState.Init(hsm)
	return hsm
}

func (h *Hsm) Transition(newState State) {
	h.currentState = newState
	newState.Init(h)
}

func (h *Hsm) Dispatch(e Event) {
	if !h.currentState.Event(h, e) {
		fmt.Println("Event not handled")
	}
}

// Define event types
const (
	EventA = iota
	EventB
	EventC
	EventTimeout
)

// Define states
type StateA struct {
	timer *time.Timer
}

func (s *StateA) Init(hsm *Hsm) {
	fmt.Println("StateA: Init")
	s.timer = time.NewTimer(2 * time.Second)
	go func() {
		<-s.timer.C
		hsm.Dispatch(EventTimeout)
	}()
}

func (s *StateA) Event(hsm *Hsm, e Event) bool {
	switch e {
	case EventA:
		fmt.Println("Ignoring additional EventA in StateA")
		return true
	case EventB:
		hsm.Transition(&StateB{})
		s.timer.Stop()
		return true
	case EventTimeout:
		fmt.Println("StateA: Timeout occurred, transitioning to StateC")
		hsm.Transition(&StateC{})
		return true
	default:
		return false
	}
}

type StateB struct{}

func (s *StateB) Init(hsm *Hsm) {
	fmt.Println("StateB: Init")
}

func (s *StateB) Event(hsm *Hsm, e Event) bool {
	switch e {
	case EventA:
		hsm.Transition(&StateA{})
		return true
	case EventC:
		fmt.Println("Ignoring EventC in StateB")
		return true
	default:
		return false
	}
}

type StateC struct{}

func (s *StateC) Init(hsm *Hsm) {
	fmt.Println("StateC: Init")
}

func (s *StateC) Event(hsm *Hsm, e Event) bool {
	switch e {
	case EventB:
		hsm.Transition(&StateB{})
		return true
	default:
		return false
	}
}

func main() {
	// Create a new state machine
	sm := NewHsm(&StateA{})

	// Send events to the state machine
	go func() {
		sm.Dispatch(EventA)
		time.Sleep(1 * time.Second)
		sm.Dispatch(EventA)
		time.Sleep(3 * time.Second)
		sm.Dispatch(EventB)
	}()

	// Wait for state machine events to stabilize
	select {
	case <-time.After(5 * time.Second):
		fmt.Println("State machine operations completed.")
	}
}
