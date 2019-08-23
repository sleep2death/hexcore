package actions

import (
	"math/rand"
	"testing"

	"github.com/sleep2death/hexcore/cards"
	"gopkg.in/go-playground/assert.v1"
)

func TestWaitForCard(t *testing.T) {
	input := make(chan string)

	hand, _ := cards.CreatePile([]string{"Strike", "Strike", "Strike", "Strike", "Strike", "Defend", "Defend", "Defend", "Defend", "Bash"})

	card, _ := hand.GetCard(0)
	cardID := card.ID()

	ctx := &Context{}
	ctx.Cards = make([]*cards.Pile, 4)
	ctx.Cards[cards.Hand] = hand
	ctx.InputC = input

	Execute(WaitForCard, ctx)
	err := <-ctx.ErrC
	assert.Equal(t, ErrInputTimeout, err)

	Execute(WaitForCard, ctx)
	input <- cardID

	//errc should be closed
	err = <-ctx.ErrC
	assert.Equal(t, nil, err)

	Execute(WaitForCard, ctx)
	input <- "ABC"

	//errc should be closed
	err = <-ctx.ErrC
	assert.Equal(t, cards.ErrCardNotExist, err)
}

func TestStartBattle(t *testing.T) {
	// input := make(chan string)

	seed := rand.New(rand.NewSource(9012))
	deck, _ := cards.CreatePile([]string{"Strike", "Strike", "Strike", "Strike", "Strike", "Defend", "Defend", "Defend", "Defend", "Bash"})

	ctx := &Context{}

	Execute(StartBattle, ctx)
	err := <-ctx.ErrC
	assert.Equal(t, ErrSeedNotSet, err)

	ctx.Seed = seed
	Execute(StartBattle, ctx)
	err = <-ctx.ErrC
	assert.Equal(t, cards.ErrPileIsNilOrEmpty, err)

	ctx.Deck = deck
	Execute(StartBattle, ctx)
	err = <-ctx.ErrC
	assert.Equal(t, nil, err)
	assert.Equal(t, "[Bash Defend Defend Strike Strike Strike Strike Defend Defend Strike]", ctx.Cards[cards.Draw].String())
}
