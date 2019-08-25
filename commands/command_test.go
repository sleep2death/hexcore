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
	select {
	case n := <-ctx.inputc:
		buf := make([]byte, 4)
		binary.LittleEndian.PutUint16(buf, uint16(n))
		ctx.outputc <- buf
		return []Command{WaitForInput}, nil
	case <-time.After(time.Second * 5):
		return nil, ErrTimeout
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
			select {
			case err = <-errc:
				break receive
			case out := <-ctx.Output():
				v := binary.LittleEndian.Uint16(out)
				if v > 20 {
					break receive
				}
				log.Printf("output is %v", binary.LittleEndian.Uint16(out))
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

	assert.Equal(t, ErrTimeout, err)
}
