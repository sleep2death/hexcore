package hexcore

import (
	"sync"
)

// State - hold all status data of the player
// it may access by different goroutines
// so keep in mind about the concurrency safe
type State struct {
	num int
	mu  sync.Mutex
}

// Num of the state
func (s *State) Num() int {
	s.mu.Lock()
	n := s.num
	s.mu.Unlock()
	return n
}

// SetNum of the state
func (s *State) SetNum(i int) {
	s.mu.Lock()
	s.num = i
	s.mu.Unlock()
}

var store = &Store{}

// GetStore of all the states
func GetStore() *Store {
	return store
}

// Store the state of the execution
type Store struct {
	mu     sync.Mutex
	idx    int
	states []*State
}

// AddState of the store
func (s *Store) AddState(state *State) int {
	s.mu.Lock()
	s.states = append(s.states, state)
	i := len(s.states) - 1
	s.mu.Unlock()
	return i
}

// State of the store
func (s *Store) State(idx int) *State {
	s.mu.Lock()
	// TODO: idx validation
	st := s.states[idx]
	s.mu.Unlock()

	return st
}
