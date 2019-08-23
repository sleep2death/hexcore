package actions

import (
	"errors"
	"time"

	"github.com/sleep2death/hexcore/cards"
)

var (
	// ErrSeedNotSet -
	ErrSeedNotSet = errors.New("seed not set")
)

// StartBattle acttion
func StartBattle(ctx *Context) (res []Action) {
	if ctx.Seed == nil {
		ctx.ErrC <- ErrSeedNotSet
		return
	}
	// init cards in battle
	ctx.Cards = make([]*cards.Pile, 4)
	for i := 0; i < 4; i++ {
		ctx.Cards[i], _ = cards.CreatePile(nil)
	}

	// copy deck cards into draw pile
	draw := ctx.Cards[cards.Draw]
	if err := draw.CopyCardsFrom(ctx.Deck); err != nil {
		ctx.ErrC <- err
		return
	}

	// shuffle draw pile
	draw.Shuffle(ctx.Seed)

	return
}

// WaitForCard action
func WaitForCard(ctx *Context) (res []Action) {
	hand := ctx.Cards[cards.Hand]
	select {
	case id := <-ctx.InputC:
		card, _, _ := hand.FindCard(id)
		if card != nil {
			// TODO: real card play handler is needed here
			ctx.Card = card
			return []Action{Trace}
		}
		ctx.ErrC <- cards.ErrCardNotExist
	case <-time.After(time.Second * 5):
		ctx.ErrC <- ErrInputTimeout
	}
	return
}
