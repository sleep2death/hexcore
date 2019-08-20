package hexcore

import (

	// "github.com/sleep2death/hexcore/actors"

	"errors"
	"log"
	"sync/atomic"
	"time"

	"github.com/sleep2death/hexcore/actions"
	"github.com/sleep2death/hexcore/cards"
)

var seed int64 = 9012

// InputTimeout -
var InputTimeout = time.Second * 10

// ErrInputTimeout -
var ErrInputTimeout = errors.New("user input timeout")

// ErrInputInvalid -
var ErrInputInvalid = errors.New("user input invalid")

// State is the type of battle state
const (
	// WaitForNone is the state of rejecting user input
	WaitForNone int32 = iota
	// WaitForPlay is the state of waiting user to select a card
	WaitForCardPlay
	// WaitForSelect is the state of waiting user to select a target
	WaitForCardSelect
)

// Battle holds all card piles
type Battle struct {
	cards      *cards.Manager
	inputState int32

	actions   chan []actions.Action
	err       chan error
	CardInput chan string
}

// CreateBattle Manager
func CreateBattle(manager *cards.Manager) *Battle {
	b := &Battle{
		cards:      manager,
		inputState: WaitForNone,
	}

	b.CardInput = make(chan string, 1)
	b.err = make(chan error, 1)
	// cards.SetActionDispacher(b.dispatcher)
	return b
}

// InputState of the battle
func (b *Battle) InputState() int32 {
	return atomic.LoadInt32(&(b.inputState))
}

// SetInputState -
func (b *Battle) SetInputState(state int32) {
	atomic.StoreInt32(&(b.inputState), state)
}

// Update the battle loop
func (b *Battle) Update() error {
	for {
		select {
		case id := <-b.CardInput:
			switch b.InputState() {
			case WaitForCardPlay:
				card := b.cards.GetCardByID(id, cards.Hand)
				log.Printf("playing card <%v>", card)
			case WaitForCardSelect:
				log.Printf("selecting card <%s>", id)
			}
		case err := <-b.err:
			log.Printf("error fired: <%v>", err)
		default:
		}
	}
}
