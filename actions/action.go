package actions

import (
	"errors"
	"time"

	"github.com/sleep2death/hexcore/cards"
)

var (
	// ErrActionListIsEmpty -
	ErrActionListIsEmpty = errors.New("action list is empty or nil")
	// ErrWaitingForUserInput -
	ErrWaitingForUserInput = errors.New("waiting for user input")
	// ErrInputTimeout -
	ErrInputTimeout = errors.New("waiting for user input timeout")
)

// Action -
type Action interface {
	// Exec - execute the action
	// @output is the channel for sending actions back to client
	// @errorc is the channel for taking care of errors in the action
	Exec() ([]Action, error)
}

// StartBattleAction - copy all cards for deck to draw, then shuffle draw
type StartBattleAction struct {
	Deck *cards.Pile
	Draw *cards.Pile
}

// Exec -
func (a *StartBattleAction) Exec() ([]Action, error) {
	a.Draw.Clear()

	err := a.Draw.CopyCardsFrom(a.Deck)
	if err != nil {
		return nil, err
	}

	// TODO: sending the action back to client
	// by using the output channel

	return []Action{
		&ShuffleAction{
			Pile: a.Draw,
		},
	}, nil
}

// StartTurnAction - draw 5 cards from draw pile into hand
type StartTurnAction struct {
	Draw    *cards.Pile
	Hand    *cards.Pile
	Discard *cards.Pile
}

// Exec -
func (a *StartTurnAction) Exec() ([]Action, error) {
	err := a.Hand.Draw(a.Draw, 5)
	if err == cards.ErrDrawNumber {
		return nil, err
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
		}, nil

		// bot discard and draw pile is empty
	} else if err == cards.ErrNotEnoughCards {
		return nil, err
	}
	// TODO: sending the action back to client
	// by using the output channel
	return nil, nil
}

const (
	// WaitForAction state
	WaitForAction int32 = iota
	// WaitForPlay state
	WaitForPlay
	// WaitForDiscard state
	WaitForDiscard
	// WaitForExaust state
	WaitForExaust
)

// WaitForPlayAction - wait for player input the card
type WaitForPlayAction struct {
	Hand  *cards.Pile
	Input chan string
}

// Exec -
func (a *WaitForPlayAction) Exec() ([]Action, error) {
	select {
	case id := <-a.Input:
		_, _, err := a.Hand.FindCard(id)
		if err != nil {
			return nil, err
		}
		return nil, nil
	case <-time.After(time.Second * 5):
		return nil, ErrInputTimeout
	}
}

// PlayCardAction -
// type PlayCardAction struct {
// 	Card *cards.Card
// }

// // Exec -
// func (a *PlayCardAction) Exec() ([]Action, error) {
// 	log.Printf("Play card -> %v", a.Card)
// 	return nil, nil
// }

// DiscardAction -
type DiscardAction struct {
	ID      string
	Hand    *cards.Pile
	Discard *cards.Pile
}

// Exec -
func (a *DiscardAction) Exec() ([]Action, error) {
	if err := a.Discard.Pick(a.Hand, a.ID); err != nil {
		return nil, err
	}
	return nil, nil
}

// DiscardAllAction -
type DiscardAllAction struct {
	Hand    *cards.Pile
	Discard *cards.Pile
}

// Exec -
func (a *DiscardAllAction) Exec() ([]Action, error) {
	num := a.Hand.Num()
	if num == 0 {
		return nil, cards.ErrNotEnoughCards
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
	return discardBatch, nil
}

// # Card & Pile Actions #

// ShuffleAction -
type ShuffleAction struct {
	Pile *cards.Pile
}

// Exec -
func (a *ShuffleAction) Exec() ([]Action, error) {
	a.Pile.Shuffle()
	return nil, nil
}

// DrawAction -
type DrawAction struct {
	Pile     *cards.Pile
	DrawFrom *cards.Pile
	N        int
}

// Exec -
func (a *DrawAction) Exec() ([]Action, error) {
	err := a.Pile.Draw(a.DrawFrom, a.N)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
