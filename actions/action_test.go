package actions

import (
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/sleep2death/hexcore/cards"
	"github.com/stretchr/testify/assert"
)

// Do the actions
func Do(actions []Action) error {
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
		},
	})
	assert.Equal(t, 5, draw.Num())
	assert.Equal(t, 5, hand.Num())
	assert.Equal(t, 0, discard.Num())

	// time.Sleep(time.Second * 6)
}

type AsyncAction interface {
	Exec(errc chan<- error) []AsyncAction
}

func DoAsync(actions []AsyncAction, errc chan<- error) {
	if actions == nil || len(actions) == 0 {
		errc <- ErrActionListIsEmpty
	}
	for _, a := range actions {
		as := a.Exec(errc)
		if as != nil {
			DoAsync(as, errc)
		}
	}
}

type WaitForInputAction struct {
	Hand  *cards.Pile
	Input chan string
}

func (a *WaitForInputAction) Exec(errc chan<- error) []AsyncAction {
	select {
	case id := <-a.Input:
		card, _, err := a.Hand.FindCard(id)
		if err != nil {
			errc <- err
			return nil
		}
		return []AsyncAction{
			&PlayCardAction{
				Card: card,
			},
		}
	case <-time.After(time.Second * 3):
		errc <- ErrInputTimeout
		return nil
	}
}

type PlayCardAction struct {
	Card *cards.Card
}

func (a *PlayCardAction) Exec(errc chan<- error) []AsyncAction {
	log.Printf("Play card -> %v", a.Card)
	return nil
}

func TestAsyncAction(t *testing.T) {
	var action AsyncAction

	seed := rand.New(rand.NewSource(9012))
	hand, _ := cards.CreatePile(seed, []string{"Strike", "Strike", "Strike", "Strike", "Strike", "Defend", "Defend", "Defend", "Defend", "Bash"})
	input := make(chan string)

	action = &WaitForInputAction{
		Hand:  hand,
		Input: input,
	}

	errc := make(chan error)

	// card, _ := hand.GetCard(0)
	// cardID := card.ID()
	go DoAsync([]AsyncAction{action}, errc)

	time.Sleep(time.Second * 1)
	input <- "ABC"
	err := <-errc
	t.Log(err.Error())

	// errc = make(chan error)
	// DoAsync(action, errc)
	// err = <-errc
	// assert.Equal(t, ErrInputTimeout, err)
}
