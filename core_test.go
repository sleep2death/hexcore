package hexcore

import (
	"encoding/binary"
	"log"
	"testing"
	"time"

	"github.com/sleep2death/hexcore/actions"

	"github.com/stretchr/testify/assert"
)

type update struct {
	delta int
}

func (a *update) Exec(ctx *actions.Context) ([]actions.Action, error) {
	state := actions.GetStore().State(ctx.ID())
	state.SetNum(state.Num() + a.delta)

	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, uint32(state.Num()))

	// log.Printf("updated num: %d:", state.Num)

	select {
	case ctx.Output() <- bs:
		// log.Printf("sending num: %d:", state.Num)
		return nil, nil
	case <-time.After(time.Second * 5): // timeout
		return nil, actions.ErrTimeout
	}
}

func TestChainWithInputCancel(t *testing.T) {
	// a test starting state
	state := &actions.State{}
	state.SetNum(5)

	errc, inputc, outputc := Start(nil, state)

	// done channel: close it to stop sender from sending data to execution
	done := make(chan struct{})

	go func() {
	receiver: // continuously reading data from execution
		for {
			select {
			case err := <-errc: // read the execution result
				assert.Equal(t, actions.ErrCanceled, err)
				assert.Equal(t, 15, state.Num())
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
	state := &actions.State{}
	state.SetNum(5)

	errc, inputc, _ := Start(&actions.WaitForInput{}, state)

	// done channel: close it to stop sender from sending data to execution
	done := make(chan struct{})

	go func() {
	receiver: // continuously reading data from execution
		for {
			select {
			case err := <-errc: // read the execution result
				assert.Equal(t, actions.ErrTimeout, err)
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
