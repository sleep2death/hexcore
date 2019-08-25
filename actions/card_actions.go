package actions

import "github.com/sleep2death/hexcore/cards"

// playAction -
var playAction = map[string]func(ctx *Context) []Action{
	"Strike": strike,
}

// PlayCard -
func PlayCard(ctx *Context) []Action {
	if ctx.input == nil {
		ctx.errc <- ErrNilInput
		return nil
	}

	card, _, err := ctx.piles[cards.Hand].FindCard(ctx.input.CardID)

	if err != nil {
		ctx.errc <- cards.ErrCardNotExist
		return nil
	}

	action := playAction[card.Name()]
	if action == nil {
		ctx.errc <- ErrActionNotFound
		return nil
	}

	return []Action{action}
}

// strike card play action
func strike(ctx *Context) []Action {
	return nil
}
