package commands

import (
	"encoding/binary"
	"log"
	"testing"
	"time"

	"gopkg.in/go-playground/assert.v1"
)

func WaitForInput(ctx *Context) ([]Command, error) {
	// log.Println("Waiting For Input")
	var buf []byte
	select {
	case n := <-ctx.inputc:
		buf = make([]byte, 4)
		binary.LittleEndian.PutUint16(buf, uint16(n))
		ctx.data = buf
		return []Command{Output}, nil
	case <-time.After(time.Second * 5):
		return nil, ErrTimeout
	case <-ctx.done:
		return nil, ErrCanceled
	}
}

func Output(ctx *Context) ([]Command, error) {
	select {
	case ctx.outputc <- ctx.data:
		return []Command{WaitForInput}, nil
	case <-ctx.done:
		return nil, ErrCanceled
	}
}

func TestChain(t *testing.T) {
	count := 0
	var err error
	done := make(chan struct{})

	ctx := NewContext()
	errc := Exec(WaitForInput, ctx)

	go func() {
	receive:
		for {
			out := <-ctx.Output()
			v := binary.LittleEndian.Uint16(out)
			log.Printf("output is %v", binary.LittleEndian.Uint16(out))

			if v >= 20 {
				ctx.Cancel()
				break receive
			}
		}

		close(done)
	}()

	// write input opration must be different goroutine of receive
send:
	for {
		count++

		select {
		case <-done:
			break send
		case ctx.Input() <- count:
		}
	}

	err = <-errc
	assert.Equal(t, ErrCanceled, err)
}
