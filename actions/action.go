package actions

import (
	"log"
	"time"

	"github.com/sleep2death/hexcore/cards"
)

// Execute the actions
func Execute(actions []Action, output chan<- Action, errorc chan<- error) []Action {
	if actions == nil || len(actions) == 0 {
		return nil
	}
	for _, a := range actions {
		if a != nil {
			as := a.Exec(output, errorc)
			Execute(as, output, errorc)
		}
	}

	return nil
}

// Action -
type Action interface {
	// Exec - execute the action
	// @output is the channel for sending actions back to client
	// @errorc is the channel for taking care of errors in the action
	Exec(output chan<- Action, errorc chan<- error) []Action
}

// StartBattleAction - copy all cards for deck to draw, then shuffle draw
type StartBattleAction struct {
	Deck *cards.Pile
	Draw *cards.Pile
}

// Exec -
func (a *StartBattleAction) Exec(output chan<- Action, errorc chan<- error) []Action {
	a.Draw.Clear()

	err := a.Draw.CopyCardsFrom(a.Deck)
	if err != nil {
		errorc <- err
		return nil
	}

	// TODO: sending the action back to client
	// by using the output channel

	return []Action{
		&ShuffleAction{
			Pile: a.Draw,
		},
	}
}

// StartTurnAction - draw 5 cards from draw pile into hand
type StartTurnAction struct {
	Draw    *cards.Pile
	Hand    *cards.Pile
	Discard *cards.Pile
}

// Exec -
func (a *StartTurnAction) Exec(output chan<- Action, errorc chan<- error) []Action {
	err := a.Hand.Draw(a.Draw, 5)
	if err == cards.ErrDrawNumber {
		errorc <- err
		return nil
	} else if err == cards.ErrNotEnoughCards && a.Discard.Num() > 0 {
		// if draw pile in not enough,
		// draw all the left cards of the draw pile into hand
		// then shuffle the discard pile and put it in the bottom of the draw pile
		left := 5 - a.Draw.Num()
		a.Hand.Draw(a.Draw, a.Draw.Num())
		return []Action{
			&ShuffleAction{
				Pile: a.Discard,
			},
			&DrawAction{
				Pile:     a.Draw,
				DrawFrom: a.Discard,
				N:        a.Discard.Num(),
			},
			&DrawAction{
				Pile:     a.Hand,
				DrawFrom: a.Draw,
				N:        left,
			},
		}

	}
	// TODO: sending the action back to client
	// by using the output channel
	return nil
}

// WaitForPlayAction - wait for player input the card
type WaitForPlayAction struct {
	Hand  *cards.Pile
	Input chan string
}

// Exec -
func (a *WaitForPlayAction) Exec(output chan<- Action, erroc chan<- error) []Action {
	select {
	case id := <-a.Input:
		log.Printf("Input Card ID: %s", id)
		return nil
	case <-time.After(time.Second * 5):
		log.Print("TimeOut")
		return nil
	}
}

// DiscardAction -
type DiscardAction struct {
	ID      string
	Hand    *cards.Pile
	Discard *cards.Pile
}

// Exec -
func (a *DiscardAction) Exec(output chan<- Action, errorc chan<- error) []Action {
	if err := a.Discard.Pick(a.Hand, a.ID); err != nil {
		errorc <- err
	}
	return nil
}

// DiscardAllAction -
type DiscardAllAction struct {
	Hand    *cards.Pile
	Discard *cards.Pile
}

// Exec -
func (a *DiscardAllAction) Exec(output chan<- Action, errorc chan<- error) []Action {
	num := a.Hand.Num()
	if num == 0 {
		errorc <- cards.ErrNotEnoughCards
		return nil
	}
	discardBatch := make([]Action, num)
	for i := 0; i < num; i++ {
		card, _ := a.Hand.GetCard(i)
		discardBatch[i] = &DiscardAction{
			Hand:    a.Hand,
			Discard: a.Discard,
			ID:      card.ID(),
		}
	}
	return discardBatch
}

// # Card & Pile Actions #

// ShuffleAction -
type ShuffleAction struct {
	Pile *cards.Pile
}

// Exec -
func (a *ShuffleAction) Exec(output chan<- Action, errorc chan<- error) []Action {
	a.Pile.Shuffle()
	return nil
}

// DrawAction -
type DrawAction struct {
	Pile     *cards.Pile
	DrawFrom *cards.Pile
	N        int
}

// Exec -
func (a *DrawAction) Exec(output chan<- Action, errorc chan<- error) []Action {
	err := a.Pile.Draw(a.DrawFrom, a.N)
	if err != nil {
		errorc <- err
	}
	return nil
}
