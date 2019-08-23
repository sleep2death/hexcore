package actions

import (
	"errors"
	"log"
	"math/rand"

	"github.com/sleep2death/hexcore/actors"
	"github.com/sleep2death/hexcore/cards"
)

var (
	// ErrActionListIsEmpty -
	ErrActionListIsEmpty = errors.New("action list is empty or nil")
	// ErrInputTimeout -
	ErrInputTimeout = errors.New("waiting for user input timeout")
)

// Data for callback
type Data interface {
	Send()
}

// Context of the action
type Context struct {
	Seed *rand.Rand

	Deck  *cards.Pile   // deck cards in player's pocket :)
	Cards []*cards.Pile // cards in the battle

	Player   *actors.Player    // player
	Monsters []*actors.Monster // monsters

	Card   *cards.Card   // current selected card
	Target *actors.Actor // current selected target

	InputC <-chan string

	ErrC chan error
	OutC chan Data
}

// Action -
type Action func(ctx *Context) []Action

// Run action chain
func run(action Action, ctx *Context) {
	next := action(ctx)

	if next == nil || len(next) == 0 {
		return
	}

	for _, n := range next {
		run(n, ctx)
	}
}

// Execute action async
func Execute(action Action, ctx *Context) {
	ctx.ErrC = make(chan error)
	ctx.OutC = make(chan Data)

	go func() {
		defer close(ctx.ErrC)
		defer close(ctx.OutC)

		run(action, ctx)
	}()
}

// Receive example:
// func Receive(errc <-chan error, outc <-chan Data) {
// 	go func() {
// 		for {
// 			select {
// 			case err := <-errc:
// 				if err != nil {
// 				}
// 				return
// 			case <-outc:
// 				//TODO: send the action result back to client
// 			default:
// 			}
// 		}
// 	}()
// }

// Trace action
// Do -
func Trace(ctx *Context) []Action {
	log.Printf("ctx: <%+v>", ctx)
	return nil
}
