package commands

import (
	"encoding/binary"
	"log"
	"testing"
	"time"

	"gopkg.in/go-playground/assert.v1"
)

func WaitForInput(ctx *Context) ([]Command, error) {
	log.Println("Waiting For Input")
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

			default:
				log.Print("receiving...")
				// make an execution or send data to inptut channel
			}
		}

		close(done)
	}()

	// write input opration must be different goroutine of receive
send:
	for {
		select {
		case <-done:
			break send
		case ctx.Input() <- count:
			count++
		default:
			log.Print("sending...")
		}
	}

	assert.Equal(t, ErrTimeout, err)
}
