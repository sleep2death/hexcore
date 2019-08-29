package hexcore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TempAction struct {
}

func (a *TempAction) Exec(ctx *Context) ([]Action, error) {
	return nil, nil
}

func TestInputAction(t *testing.T) {
	in := make(chan Action)
	out := make(chan []byte)
	id := 1

	ctx := NewContext(in, out, id)
	action := &WaitForInput{}

	// timeout -
	_, err := action.Exec(ctx)
	assert.Equal(t, ErrTimeout, err)

	// send action to input channel
	nextAction := &TempAction{}
	go func() {
		in <- nextAction
	}()

	actions, err := action.Exec(ctx)
	assert.Equal(t, nextAction, actions[0])
	assert.Equal(t, nil, err)

	// close input channel outside
	go func() {
		close(in)
	}()

	_, err = action.Exec(ctx)
	assert.Equal(t, ErrCanceled, err)
}
