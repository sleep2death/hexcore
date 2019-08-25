package actions

import (
	"time"

	"github.com/sleep2death/hexcore/cards"
)

// start battle -> start turn -> wait for player -> play card(s) or drink potion(s) -> turn end

// StartBattle action
func StartBattle(ctx *Context) (res []Action) {
	// init cards in battle
	ctx.mux.Lock()
	for i := 1; i < 5; i++ {
		ctx.piles[i], _ = cards.CreatePile(nil)
	}

	// copy deck cards into draw pile
	draw := ctx.piles[cards.Draw]
	deck := ctx.piles[cards.Deck]

	ctx.mux.Unlock()

	if err := draw.CopyCardsFrom(deck); err != nil {
		ctx.errc <- err
		return
	}

	// shuffle draw pile
	draw.Shuffle(ctx.seed)

	// TODO: battle start triggers need to handle here

	res = []Action{StartTurn}
	return
}

// StartTurn action
func StartTurn(ctx *Context) (res []Action) {
	draw := ctx.piles[cards.Draw]
	hand := ctx.piles[cards.Hand]
	discard := ctx.piles[cards.Discard]

	err := hand.Draw(draw, 5)
	// if err != nil, err must be cards.ErrNotEnoughCards
	if err != nil && discard.Num() > 0 {
		left := 5 - draw.Num()
		// draw all left cards of the draw pile
		hand.Draw(draw, draw.Num())
		// shuffle discard pile
		discard.Shuffle(ctx.seed)
		// move discard to draw
		draw.Draw(discard, discard.Num())

		// if left-to-draw num is larger than draw pile num, then draw all of the draw pile
		if left > draw.Num() {
			left = draw.Num()
		}

		// draw the left
		hand.Draw(draw, left)
		return
	}

	// TODO: turn start triggers need to handle here
	res = []Action{WaitForPlay}
	return
}

// WaitForPlay action
func WaitForPlay(ctx *Context) (res []Action) {
	select {
	case input := <-ctx.inputc:
		if input == nil {
			ctx.errc <- ErrNilInput
		}
		ctx.input = input
		return []Action{PlayCard}
	case <-time.After(time.Second * 5):
		ctx.errc <- ErrInputTimeout
	}
	return
}
