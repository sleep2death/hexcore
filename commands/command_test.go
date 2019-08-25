package commands

import (
	"errors"
	"log"
	"testing"
	"time"

	"gopkg.in/go-playground/assert.v1"
)

func WaitForInput(ctx *Context) ([]Command, error) {
	log.Println("Waiting...")
	select {
	case n := <-ctx.inputc:
		log.Printf("input is: %d", n)
		ctx.outputc <- n
		return []Command{WaitForInput}, nil
	case <-time.After(time.Second * 5):
		return nil, errors.New("execution timeout")
	}
}

func TestChain(t *testing.T) {
	count := 5
	var err error

	ctx := NewContext()
	errc := Exec(WaitForInput, ctx)

receive:
	for {
		// receive socket packet here, and send it to ctx.Input
		if count > 0 {
			ctx.Input() <- 5
			count --
			break send
		}
		select {
		// nil channel will block, so it will select the default
		// closed channel will not block, so it should set to nil
		case err = <-errc:
			if err != nil {
				break loop
			}
			errc = nil
		case out := <-ctx.Output():
			log.Printf("output is %d", out)
		default:
			// make an execution or send data to inptut channel
		}
	}

	assert.Equal(t, "execution timeout", err.Error())
}
