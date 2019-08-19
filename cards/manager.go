package cards

import (
	"errors"
	"fmt"
)

var seed int64 = 9012

// ErrDrawNumber -
var ErrDrawNumber error = errors.New("draw number should be larger than 0")

// ErrDrawIndex -
var ErrDrawIndex error = errors.New("draw index should be larger than 0 and less than len(cards) - 1")

// ErrNotEnoughCards -
var ErrNotEnoughCards error = errors.New("not enough card(s) to draw")

// ErrCardNotExist -
var ErrCardNotExist error = errors.New("card doesn't exist in the pile")

// ErrPileIsNil -
var ErrPileIsNil error = errors.New("target pile is nil")

// Manager of all the cards in player's pocket
type Manager struct {
	deck    *Pile
	draw    *Pile
	hand    *Pile
	discard *Pile
	exaust  *Pile
}

// Create the card manager by filling the deck with given cards
func (m *Manager) Create(cardSet []string) error {
	cap := len(cardSet)
	dCap := cap * 4
	m.deck = &Pile{seed: seed, cards: make([]Card, 0, cap)}
	m.draw = &Pile{seed: seed, cards: make([]Card, 0, dCap)}
	m.hand = &Pile{seed: seed, cards: make([]Card, 0, dCap)}
	m.discard = &Pile{seed: seed, cards: make([]Card, 0, dCap)}
	m.exaust = &Pile{seed: seed, cards: make([]Card, 0, dCap)}
	return m.deck.CreateCardByName(cardSet)
}

func (m *Manager) String() string {
	return fmt.Sprintf("CardManager ->\ndeck: %v\ndraw: %v\nhand: %v\ndiscard: %v\nexaust: %v", m.deck.cards, m.draw.cards, m.hand.cards, m.discard.cards, m.exaust.cards)
}

// Shuffle - copy all cards from the deck pile to draw pile
func (m *Manager) Shuffle() error {
	if m.deck == nil || m.deck.CardsNum() <= 0 {
		return ErrPileIsNil
	}

	m.draw = &Pile{
		seed:  m.deck.seed,
		cards: make([]Card, len(m.deck.cards)),
	}

	for i, card := range m.deck.cards {
		m.draw.cards[i] = card.Copy()
	}

	m.draw.Shuffle()
	return nil
}

// ReShuffle all cards in the discard pile, then draw all the cards into draw pile
func (m *Manager) ReShuffle() error {
	m.discard.Shuffle()
	return m.discard.Draw(m.discard.CardsNum(), m.draw)
}

// Draw n cards from draw pile into hand pile
func (m *Manager) Draw(n int) error {
	// TODO: hand pile may have some cards number limit
	return m.draw.Draw(n, m.hand)
}

// Exaust cards from hand pile to exaust pile
func (m *Manager) Exaust(card Card) error {
	idx := m.hand.FindCardByID(card.ID())
	if idx < 0 {
		return ErrCardNotExist
	}
	// TODO: exaust action triger here
	// draw the card from hand pile into exaust pile
	return m.hand.DrawCard(idx, m.exaust)
}

// ExaustAll cards from hand
func (m *Manager) ExaustAll() error {
	// should exaust the card one by one
	// some cards will trigger actions when exausted
	for m.hand.CardsNum() > 0 {
		if err := m.Exaust(m.hand.cards[0]); err != nil {
			return err
		}
	}

	return nil
}

// Discard cards from hand pile to discard pile
func (m *Manager) Discard(card Card) error {
	idx := m.hand.FindCardByID(card.ID())
	if idx < 0 {
		return ErrCardNotExist
	}
	// TODO: discard action trigger here
	// draw the card from hand pile into exaust pile
	return m.hand.DrawCard(idx, m.discard)
}

// DiscardAll cards from hand
func (m *Manager) DiscardAll() error {
	// should discard the card one by one
	// some cards will trigger actions when discarded
	for m.hand.CardsNum() > 0 {
		if err := m.Discard(m.hand.cards[0]); err != nil {
			return err
		}
	}

	return nil
}
