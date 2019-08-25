package actions

import (
	"testing"
	"time"

	"github.com/sleep2death/hexcore/cards"
	"gopkg.in/go-playground/assert.v1"
)

func TestStartBattle(t *testing.T) {
	piles := make([]*cards.Pile, 4)
	ctx, err := NewContext(int64(9012), piles)
	assert.Equal(t, ErrInitCards, err)

	deck, _ := cards.CreatePile([]string{"Strike", "Strike", "Strike", "Strike", "Strike", "Defend", "Defend", "Defend", "Defend", "Bash"})

	piles = make([]*cards.Pile, 5)
	piles[cards.Deck] = deck

	// input wrong card id
	ctx, _ = NewContext(int64(9012), piles)
	inputc, errc, _ := Execute(StartBattle, ctx)
	// Must wait a moment, for turn start function return
	time.Sleep(time.Millisecond * 10)
	inputc <- &Input{CardID: "ABC"}

	err = <-errc
	assert.Equal(t, cards.ErrCardNotExist, err)
	assert.Equal(t, "[Bash Defend Defend Strike Strike]", ctx.piles[cards.Draw].String())
	assert.Equal(t, "[Strike Strike Defend Defend Strike]", ctx.piles[cards.Hand].String())

	// input timeout
	ctx, _ = NewContext(int64(9012), piles)
	_, errc, _ = Execute(StartBattle, ctx)
	err = <-errc
	assert.Equal(t, ErrInputTimeout, err)
	assert.Equal(t, "[Bash Defend Defend Strike Strike]", ctx.piles[cards.Draw].String())
	assert.Equal(t, "[Strike Strike Defend Defend Strike]", ctx.piles[cards.Hand].String())

	// execute actions in the loop
	errc = nil
	for {
		// wait for previous errc closing, or return
		if errc != nil {
			err = <-errc
			assert.Equal(t, nil, err)
			assert.Equal(t, "[Bash Defend Defend Strike Strike]", ctx.piles[cards.Draw].String())
			assert.Equal(t, "[Strike Strike Defend Defend Strike]", ctx.piles[cards.Hand].String())

			errc = nil

			break
		}
		ctx, _ = NewContext(int64(9012), piles)
		inputc, errc, _ = Execute(StartBattle, ctx)

		card := ctx.CardByIndex(cards.Hand, 0)
		inputc <- &Input{CardID: card.ID()}
	}

}

func TestPlayCard(t *testing.T) {
	piles := make([]*cards.Pile, 4)
	ctx, err := NewContext(int64(9012), piles)
	assert.Equal(t, ErrInitCards, err)

	// create the default deck
	deck, _ := cards.CreatePile([]string{"Strike", "Strike", "Strike", "Strike", "Strike", "Defend", "Defend", "Defend", "Defend", "Bash"})

	piles = make([]*cards.Pile, 5)
	piles[cards.Deck] = deck
	// Play the first card
	ctx, _ = NewContext(int64(9012), piles)
	inputc, errc, _ := Execute(StartBattle, ctx)

	// Must wait a moment, for turn start function return
	time.Sleep(time.Millisecond * 10)

	// Now current waiting action is WaitForPlay
	// Get the first card of hand by locking/unlocking the mutex
	card := ctx.CardByIndex(cards.Hand, 0)
	inputc <- &Input{CardID: card.ID()}

	err = <-errc
	assert.Equal(t, ErrActionNotFound, err)
}
