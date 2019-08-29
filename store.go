package hexcore

import (
	"math/rand"
	"sync"
)

// PileName -
type PileName int

const (
	// Deck pile
	Deck PileName = iota
	// Draw pile
	Draw
	// Hand pile
	Hand
	// Discard pile
	Discard
	// Exhaust pile
	Exhaust
)

// State - hold all status data of the player
// it may access by different goroutines
// so keep in mind about the concurrency safe
type State struct {
	mu  sync.Mutex
	num int

	deck    *Pile
	draw    *Pile
	hand    *Pile
	discard *Pile
	exhaust *Pile
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

// SetPile of the state
func (s *State) SetPile(name PileName, pile *Pile) {
	s.mu.Lock()
	switch name {
	case Deck:
		s.deck = pile
	case Draw:
		s.draw = pile
	case Hand:
		s.hand = pile
	case Discard:
		s.discard = pile
	case Exhaust:
		s.exhaust = pile
	}
	s.mu.Unlock()
}

// GetPile of the state
func (s *State) GetPile(name PileName) (pile *Pile) {
	s.mu.Lock()
	switch name {
	case Deck:
		pile = s.deck
	case Draw:
		pile = s.draw
	case Hand:
		pile = s.hand
	case Discard:
		pile = s.discard
	case Exhaust:
		pile = s.exhaust
	}
	s.mu.Unlock()
	return pile
}

// Shuffle the pile of the state
func (s *State) Shuffle(name PileName, seed *rand.Rand) {
	pile := s.GetPile(name)

	s.mu.Lock()
	pile.Shuffle(seed)
	s.mu.Unlock()
}

// Draw the card from one pile to another
func (s *State) Draw(from PileName, to PileName) (Card, error) {
	pFrom := s.GetPile(from)
	pTo := s.GetPile(to)

	s.mu.Lock()
	card, error := pTo.Draw(pFrom)
	s.mu.Unlock()
	return card, error
}

// Pick the card from one pile to another
func (s *State) Pick(id string, from PileName, to PileName) (Card, error) {
	pFrom := s.GetPile(from)
	pTo := s.GetPile(to)

	s.mu.Lock()
	card, error := pTo.Pick(id, pFrom)
	s.mu.Unlock()
	return card, error
}

// Copy one pile to another
func (s *State) Copy(from PileName, to PileName) {
	pFrom := s.GetPile(from)
	s.SetPile(to, pFrom.Copy())
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
