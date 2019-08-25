package actions

import (
	"errors"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/sleep2death/hexcore/actors"
	"github.com/sleep2death/hexcore/cards"
)

var (
	// ErrInitCards -
	ErrInitCards = errors.New("initial cards of context are illigal")
	// ErrActionListIsEmpty -
	ErrActionListIsEmpty = errors.New("action list is empty or nil")
	// ErrInputTimeout -
	ErrInputTimeout = errors.New("waiting for user input timeout")
	// ErrSeedNotSet -
	ErrSeedNotSet = errors.New("seed not set")
	// ErrActionNotFound -
	ErrActionNotFound = errors.New("action not found")
	// ErrNilCard -
	ErrNilCard = errors.New("card is nil")
	// ErrNilInput -
	ErrNilInput = errors.New("input is nil")
)

// Data for callback
type Data interface {
	Send()
}

// Input -
type Input struct {
	// CardID -
	CardID string
	// TargetID -
	TargetID string
	// EndTurn
	EndTurn string
}

// Context of the action
type Context struct {
	seed *rand.Rand
	mux  *sync.Mutex

	piles []*cards.Pile // cards in the battle

	player  *actors.Player    // player
	monster []*actors.Monster // monsters

	input  *Input
	inputc <-chan *Input

	errc chan<- error // error output channel
	outc chan<- Data  // data output channel
}

// CardByIndex debug and unit test only
func (ctx *Context) CardByIndex(pilename cards.PileName, idx int) *cards.Card {
	ctx.mux.Lock()
	pile := ctx.piles[pilename]
	card, _ := pile.GetCard(idx)
	ctx.mux.Unlock()
	return card
}

// NewContext -
func NewContext(seed int64, piles []*cards.Pile) (*Context, error) {
	// validate cards
	if len(piles) != 5 || piles[cards.Deck] == nil || piles[cards.Deck].Num() == 0 {
		return nil, ErrInitCards
	}

	ctx := &Context{
		seed:  rand.New(rand.NewSource(seed)),
		mux:   &sync.Mutex{},
		piles: piles,
	}

	return ctx, nil
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
func Execute(action Action, ctx *Context) (chan<- *Input, <-chan error, <-chan Data) {
	inputc := make(chan *Input)
	errc := make(chan error)
	outc := make(chan Data)

	ctx.outc = outc

	ctx.input = nil

	ctx.inputc = inputc
	ctx.errc = errc
	ctx.outc = outc

	go func() {
		defer func() {
			close(errc)
			close(outc)
			close(inputc)
		}()

		run(action, ctx)
	}()

	// must wait for a moment
	// let the execution reach the first input wait point
	time.Sleep(time.Millisecond * 10)
	return inputc, errc, outc
}

// Trace action, print all the fields of the context
func Trace(ctx *Context) []Action {
	log.Printf("ctx: <%+v>", ctx)
	return nil
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
