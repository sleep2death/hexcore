package actions

import (
	"errors"
	"time"
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

// NewContext -
func NewContext(input <-chan Action, output chan<- []byte, id int) *Context {
	return &Context{
		inc:  input,
		outc: output,
		id:   id,
	}
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

// Timeout duration
var timeout = time.Second * 5

// WaitForInput action
// when previous action return is nil,
// this action will be automatically added into the execution chain
type WaitForInput struct {
}

// Exec -
func (a *WaitForInput) Exec(ctx *Context) ([]Action, error) {
	select {
	case action := <-ctx.Input():
		if action == nil {
			return nil, ErrCanceled
		}
		return []Action{action}, nil
	case <-time.After(timeout): // timeout
		return nil, ErrTimeout
	}
}

// OutputString -
type OutputString struct {
	Message string
}

// Exec -
func (a *OutputString) Exec(ctx *Context) ([]Action, error) {
	select {
	case ctx.Output() <- []byte(a.Message):
		return nil, nil
	case <-time.After(timeout): // timeout
		return nil, ErrTimeout
	}
}

// // Exec -
// func (a *PlayCard) Exec(ctx *Context) ([]Action, error) {
// 	state := store.GetStore().State(ctx.ID())
// 	card, _, _ := state.GetPile(store.Hand).FindCard(a.ID)
// 	action := GetActionByCardName(card.Name())
// }
