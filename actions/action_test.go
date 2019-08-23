package actions

import (
	"errors"
	"log"
	"math/rand"
	"sync/atomic"
	"testing"

	"github.com/sleep2death/hexcore/cards"
	"github.com/stretchr/testify/assert"
)

var state int32
var input chan WaitForCardResult

// Do the actions
func Do(actions []Action) error {
	if s := atomic.LoadInt32(&state); s > 0 {
		return errors.New("wait for input")
	}

	if actions == nil || len(actions) == 0 {
		return ErrActionListIsEmpty
	}
	for _, a := range actions {
		as, err := a.Exec()
		if err != nil {
			return err
		}
		Do(as)
	}

	return nil
}

func Receive(done <-chan struct{}, input <-chan WaitForCardResult) <-chan error {
	errc := make(chan error, 1)
	var err error
	go func() {
		select {
		case res := <-recv:
			if res.Err != nil {
				err = res.Err
			} else {
				log.Println("card result received", res)
			}
		case <-done:
		}
	}()
	// No select needed here, since errc is buffered.
	errc <- err
	return errc
}

func TestStartBattleAction(t *testing.T) {
	var action Action
	seed := rand.New(rand.NewSource(9012))

	deck, _ := cards.CreatePile(seed, []string{"Strike", "Strike", "Strike", "Strike", "Strike", "Defend", "Defend", "Defend", "Defend", "Bash"})
	draw, _ := cards.CreatePile(seed, nil)

	action = &StartBattleAction{
		Deck: deck,
		Draw: draw,
	}

	err := Do([]Action{action})
	assert.Equal(t, nil, err)
	assert.Equal(t, "[Bash Defend Defend Strike Strike Strike Strike Defend Defend Strike]", draw.String())

	action = &StartBattleAction{
		Deck: nil,
		Draw: draw,
	}

	err = Do([]Action{action})
	assert.Equal(t, cards.ErrPileIsNilOrEmpty, err)

	hand, _ := cards.CreatePile(seed, nil)
	discard, _ := cards.CreatePile(seed, nil)

	action = &StartTurnAction{
		Draw:    draw,
		Hand:    hand,
		Discard: discard,
	}
	err = Do([]Action{action})
	assert.Equal(t, cards.ErrNotEnoughCards, err)

	action = &StartBattleAction{
		Deck: deck,
		Draw: draw,
	}
	err = Do([]Action{action})
	assert.Equal(t, nil, err)
	assert.Equal(t, "[Strike Strike Strike Strike Defend Strike Defend Bash Defend Defend]", draw.String())

	action = &StartTurnAction{
		Draw:    draw,
		Hand:    hand,
		Discard: discard,
	}
	err = Do([]Action{action})
	assert.Equal(t, nil, err)
	assert.Equal(t, "[Strike Strike Strike Strike Defend]", draw.String())
	assert.Equal(t, "[Strike Defend Bash Defend Defend]", hand.String())

	action = &DiscardAllAction{
		Hand:    hand,
		Discard: discard,
	}
	err = Do([]Action{action})
	assert.Equal(t, nil, err)
	assert.Equal(t, 0, hand.Num())
	assert.Equal(t, 5, discard.Num())

	actions := []Action{
		&DrawAction{
			Pile:     hand,
			DrawFrom: draw,
			N:        3,
		},
		&DiscardAllAction{
			Hand:    hand,
			Discard: discard,
		},
	}
	err = Do(actions)
	assert.Equal(t, nil, err)
	assert.Equal(t, 2, draw.Num())
	assert.Equal(t, 0, hand.Num())

	action = &StartTurnAction{
		Draw:    draw,
		Hand:    hand,
		Discard: discard,
	}

	input := make(chan string)
	err = Do([]Action{
		action,
		&WaitForPlayAction{
			Hand:  hand,
			Input: input,
			State: &state,
		},
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, 5, draw.Num())
	assert.Equal(t, 5, hand.Num())
	assert.Equal(t, 0, discard.Num())

	// time.Sleep(time.Second * 6)
}
