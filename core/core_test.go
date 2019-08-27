package core

import (
	"encoding/binary"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type waitForInput struct {
}

func (a *waitForInput) Exec(ctx *Context) ([]Action, error) {
	select {
	case action := <-ctx.Input():
		if action == nil {
			return nil, ErrCanceled
		}
		return []Action{action}, nil
	case <-time.After(time.Second * 5): // timeout
		return nil, ErrTimeout
	}
}

type update struct {
	delta int
}

func (a *update) Exec(ctx *Context) ([]Action, error) {
	state := store.State(ctx.ID())
	state.Num += a.delta

	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, uint32(state.Num))

	// log.Printf("updated num: %d:", state.Num)

	select {
	case ctx.Output() <- bs:
		// log.Printf("sending num: %d:", state.Num)
		return []Action{&waitForInput{}}, nil
	case <-time.After(time.Second * 5): // timeout
		return nil, ErrTimeout
	}
}

func TestChainWithInputCancel(t *testing.T) {
	// a test starting state
	state := &State{Num: 5}

	errc, inputc, outputc := Start(&waitForInput{}, state)

	// done channel: close it to stop sender from sending data to execution
	done := make(chan struct{})

	go func() {
	receiver: // continuously reading data from execution
		for {
			select {
			case err := <-errc: // read the execution result
				assert.Equal(t, ErrCanceled, err)
				close(done)    // stop sender, if execution returned
				break receiver // stop receiver loop
			case data := <-outputc:
				n := binary.LittleEndian.Uint32(data)
				log.Printf("output data: %d", n)
			}
		}
	}()

sender: // continuously sending data to execution
	for i := 0; i < 5; i++ {
		select {
		case inputc <- &update{delta: 2}: // send a test update action to execution
		case <-done: // when done closed, it will break sender loop
			break sender
		}
	}

	close(inputc)
}

func TestChainWithTimeout(t *testing.T) {
	// a test starting state
	state := &State{Num: 5}

	errc, inputc, _ := Start(&waitForInput{}, state)

	// done channel: close it to stop sender from sending data to execution
	done := make(chan struct{})

	go func() {
	receiver: // continuously reading data from execution
		for {
			select {
			case err := <-errc: // read the execution result
				assert.Equal(t, ErrTimeout, err)
				close(done)    // stop sender, if execution returned
				break receiver // stop receiver loop

				// stop reading the output channel from execution,
				// it will make "update" action blocked, then timeout

				// case data := <-outputc:
				// 	n := binary.LittleEndian.Uint32(data)
				// 	log.Printf("output data: %d", n)
			}
		}
	}()

sender: // continuously sending data to execution
	for i := 0; i < 5; i++ {
		select {
		case inputc <- &update{delta: 2}: // send a test update action to execution
		case <-done: // when done closed, it will break sender loop
			break sender
		}
	}

	close(inputc)
}