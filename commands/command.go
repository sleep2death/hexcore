package commands

import (
	"errors"
	"sync"
)

var (
	// ErrTimeout -
	ErrTimeout = errors.New("execution timeout")
)

// Context of the execution
type Context struct {
	// lock
	mu sync.Mutex
	// outputc - output channel of the execution
	outputc chan []byte
	// inputc - input channel of the execution
	inputc chan int
}

// Output -
func (c *Context) Output() <-chan []byte {
	c.mu.Lock()
	o := c.outputc
	c.mu.Unlock()
	return o
}

// Input -
func (c *Context) Input() chan<- int {
	c.mu.Lock()
	i := c.inputc
	c.mu.Unlock()
	return i
}

// Close -
func (c *Context) Close() {
	c.mu.Lock()
	close(c.inputc)
	close(c.outputc)
	c.mu.Unlock()
}

// NewContext -
func NewContext() *Context {
	return &Context{
		outputc: make(chan []byte),
		inputc:  make(chan int),
	}
}

// Command function
type Command func(ctx *Context) ([]Command, error)

// exec -
func exec(comm Command, ctx *Context) error {
	next, err := comm(ctx)

	if err != nil {
		return err
	}

	if next != nil && len(next) > 0 {
		for _, n := range next {
			err = exec(n, ctx)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Exec -
func Exec(comm Command, ctx *Context) <-chan error {
	errc := make(chan error)

	go func() {
		defer close(errc)
		defer ctx.Close()

		err := exec(comm, ctx)
		errc <- err
	}()

	return errc
}
