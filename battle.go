package hexcore

import (

	// "github.com/sleep2death/hexcore/actors"

	"errors"
	"log"
	"math/rand"
	"time"

	"github.com/sleep2death/hexcore/actions"
	"github.com/sleep2death/hexcore/cards"
)

// InputTimeout -
var InputTimeout = time.Second * 10

// ErrInputTimeout -
var ErrInputTimeout = errors.New("user input timeout")

// ErrInputInvalid -
var ErrInputInvalid = errors.New("user input invalid")

// ErrCardNotFound -
var ErrCardNotFound = errors.New("user input card not found")

// PileType -
type PileType int

const (
	// DeckPile -
	DeckPile PileType = iota
	// DrawPile -
	DrawPile
	// HandPile -
	HandPile
	// DiscardPile -
	DiscardPile
	// ExaustPile -
	ExaustPile
)

// Battle holds all card piles
type Battle struct {
	seed  *rand.Rand
	cards []*cards.Pile

	input  chan string // user input channel
	errorc chan error  // error channel
}

// Exec the action
func (b *Battle) Exec(actions []actions.Action) []actions.Action {
	if actions == nil || len(actions) == 0 {
		return nil
	}

	for _, a := range actions {
		if a != nil {
			as := a.Exec(b.errorc)
			b.Exec(as)
		}
	}

	return nil
}

// Update the battle
func (b *Battle) Update() error {
	for {
		select {
		case err := <-b.errorc:
			log.Printf("action error received: %v", err)
			return err
		default:
		}
	}
}

// CreateBattle with given seed and deck pile
func CreateBattle(deck *cards.Pile) (b *Battle) {

	draw, _ := cards.CreatePile(deck.Seed, nil)
	hand, _ := cards.CreatePile(deck.Seed, nil)
	discard, _ := cards.CreatePile(deck.Seed, nil)
	exaust, _ := cards.CreatePile(deck.Seed, nil)

	b = &Battle{
		seed: deck.Seed,
		cards: []*cards.Pile{
			deck,
			draw,
			hand,
			discard,
			exaust,
		},
		input:  make(chan string),
		errorc: make(chan error),
	}
	return
}
