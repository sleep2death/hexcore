package hexcore

import (

	// "github.com/sleep2death/hexcore/actors"

	"errors"
	"fmt"
	"time"

	"github.com/sleep2death/hexcore/actions"
	"github.com/sleep2death/hexcore/actors"
	"github.com/sleep2death/hexcore/cards"
)

var seed int64 = 9012

// InputTimeout -
var InputTimeout = time.Second * 10

// ErrInputTimeout -
var ErrInputTimeout = errors.New("user input timeout")

// Battle holds all card piles
type Battle struct {
	cards *cards.Manager

	dispatcher chan actions.Action

	selectTarget chan string
	selectCards  chan []string
}

// CreateBattle Manager
func CreateBattle(cards *cards.Manager) *Battle {
	b := &Battle{
		cards: cards,

		dispatcher:   make(chan actions.Action),
		selectTarget: make(chan string),
		selectCards:  make(chan []string),
	}

	// cards.SetActionDispacher(b.dispatcher)

	return b
}

// WaitForTarget -
func (b *Battle) WaitForTarget() (c actors.Actor, err error) {
	select {
	case target := <-b.selectTarget:
		fmt.Printf("target is %s", target)
	case <-time.After(InputTimeout):
	}
	return c, ErrInputTimeout
}

// WaitForCards -
func (b *Battle) WaitForCards() (cards []cards.Card, err error) {
	select {
	case target := <-b.selectCards:
		fmt.Printf("target is %s", target)
	case <-time.After(time.Second * 5):
	}
	return cards, ErrInputTimeout
}
