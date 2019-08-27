package core

import (
	"errors"
	"sync"
)

var (
	// ErrTimeout -
	ErrTimeout = errors.New("input timeout")
	//ErrCanceled -
	ErrCanceled = errors.New("input canceled")
	//ErrNilAction -
	ErrNilAction = errors.New("input with a nil action")
)

// based on the article here: https://go101.org/article/channel-closing.html
// only sender should close the channel, never the receiver

// Context chain action channel context
// because the action is executed one by one,
// so mutex lock is not necessary
type Context struct {
	// input channel, it will be closed outside the execution
	// it's read only, so it can't be closed
	inc <-chan Action
	// output channel, it will be closed automatically,
	// when execution returned, so don't closed it in action
	outc chan<- []byte
	// context id, which can be used for finding the certain store
	id int
}

// Input channel
func (c *Context) Input() <-chan Action {
	return c.inc
}

// Output channel
func (c *Context) Output() chan<- []byte {
	return c.outc
}

// ID of the store
func (c *Context) ID() int {
	return c.id
}

// Action -
type Action interface {
	Exec(ctx *Context) ([]Action, error)
}

// Start the chain actions
func Start(action Action, state *State) (<-chan error, chan<- Action, <-chan []byte) {
	// an error channel for execution error handling
	errc := make(chan error)
	// a []byte channel for some action result data send back
	outc := make(chan []byte)
	// an input channel for executing next action
	inc := make(chan Action)

	id := store.SetState(state)

	ctx := &Context{
		outc: outc,
		inc:  inc,
		id:   id,
	}

	id++

	go func() {
		defer close(errc)
		defer close(outc)

		// execute the first action,
		// and send the last error to error channel
		err := exec(ctx, action)
		errc <- err
	}()

	return errc, inc, outc
}

// chain action execution
func exec(ctx *Context, action Action) error {
	// TODO: context and action validation
	next, err := action.Exec(ctx)
	if err != nil {
		return err
	}

	for _, action := range next {
		err = exec(ctx, action)
		// if the error is not nil, break the loop and return
		if err != nil {
			return err
		}
	}

	return nil
}

// State -
type State struct {
	Num int
}

var store = &Store{}

// Store the state of the execution
type Store struct {
	mu     sync.Mutex
	idx    int
	states []*State
}

// SetState of the store
func (s *Store) SetState(state *State) int {
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
