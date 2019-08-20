package cards

import (
	"errors"
	"fmt"
	"sync"
)

var seed int64 = 9012

// ErrDrawNumber -
var ErrDrawNumber = errors.New("draw number should be larger than 0")

// ErrDrawIndex -
var ErrDrawIndex = errors.New("draw index should be larger than 0 and less than len(cards) - 1")

// ErrNotEnoughCards -
var ErrNotEnoughCards = errors.New("not enough card(s) to draw")

// ErrCardNotExist -
var ErrCardNotExist = errors.New("card doesn't exist in the pile")

// ErrPileIsNil -
var ErrPileIsNil = errors.New("target pile is nil")

// CardPile type
type CardPile uint

const (
	// Deck Pile -
	Deck CardPile = iota
	// Draw -
	Draw
	// Hand -
	Hand
	// Discard -
	Discard
	// Exaust -
	Exaust
)

// Manager of all the cards in player's pocket
type Manager struct {
	deck    *Pile
	draw    *Pile
	hand    *Pile
	discard *Pile
	exaust  *Pile

	mux sync.Mutex
}

// CreateManager the card manager by filling the deck with given cards
func CreateManager(cardSet []string) (*Manager, error) {
	cap := len(cardSet)
	dCap := cap * 4
	m := &Manager{
		deck:    &Pile{seed: seed, cards: make([]Card, 0, cap)},
		draw:    &Pile{seed: seed, cards: make([]Card, 0, dCap)},
		hand:    &Pile{seed: seed, cards: make([]Card, 0, dCap)},
		discard: &Pile{seed: seed, cards: make([]Card, 0, dCap)},
		exaust:  &Pile{seed: seed, cards: make([]Card, 0, dCap)},
	}
	if err := m.deck.CreateCardByName(cardSet); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *Manager) String() string {
	return fmt.Sprintf("CardManager ->\ndeck: %v\ndraw: %v\nhand: %v\ndiscard: %v\nexaust: %v", m.deck.cards, m.draw.cards, m.hand.cards, m.discard.cards, m.exaust.cards)
}

// Shuffle - copy all cards from the deck pile to draw pile
func (m *Manager) Shuffle() error {
	m.mux.Lock()
	defer m.mux.Unlock()

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
	m.mux.Lock()
	defer m.mux.Unlock()

	m.discard.Shuffle()
	return m.discard.Draw(m.discard.CardsNum(), m.draw)
}

// Draw n cards from draw pile into hand pile
func (m *Manager) Draw(n int) error {
	m.mux.Lock()
	defer m.mux.Unlock()

	// TODO: hand pile may have some cards number limit
	return m.draw.Draw(n, m.hand)
}

// Exaust cards from hand pile to exaust pile
func (m *Manager) Exaust(card Card) error {
	m.mux.Lock()
	defer m.mux.Unlock()

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
	m.mux.Lock()
	defer m.mux.Unlock()
	// should exaust the card one by one
	// some cards will trigger actions when exausted
	for m.hand.CardsNum() > 0 {
		idx := m.hand.FindCardByID(m.hand.cards[0].ID())
		if idx < 0 {
			return ErrCardNotExist
		}
		if err := m.hand.DrawCard(idx, m.exaust); err != nil {
			return err
		}
	}

	return nil
}

// Discard cards from hand pile to discard pile
func (m *Manager) Discard(card Card) error {
	m.mux.Lock()
	defer m.mux.Unlock()

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
	m.mux.Lock()
	defer m.mux.Unlock()
	// should discard the card one by one
	// some cards will trigger actions when discarded
	for m.hand.CardsNum() > 0 {
		idx := m.hand.FindCardByID(m.hand.cards[0].ID())
		if idx < 0 {
			return ErrCardNotExist
		}
		if err := m.hand.DrawCard(idx, m.discard); err != nil {
			return err
		}
	}

	return nil
}

// GetCardByID -
func (m *Manager) GetCardByID(id string, pile CardPile) Card {
	var p *Pile
	switch pile {
	case Deck:
		p = m.deck
	case Draw:
		p = m.draw
	case Hand:
		p = m.hand
	case Discard:
		p = m.discard
	case Exaust:
		p = m.exaust
	default:
		return nil
	}
	m.mux.Lock()
	defer m.mux.Unlock()
	if idx := p.FindCardByID(id); idx >= 0 {
		return p.cards[idx]
	}

	return nil
}

// GetCardByIndex -
func (m *Manager) GetCardByIndex(idx int, pile CardPile) Card {
	var p *Pile
	switch pile {
	case Deck:
		p = m.deck
	case Draw:
		p = m.draw
	case Hand:
		p = m.hand
	case Discard:
		p = m.discard
	case Exaust:
		p = m.exaust
	default:
		return nil
	}
	m.mux.Lock()
	defer m.mux.Unlock()

	if idx >= 0 && idx < len(p.cards)-1 {
		return p.cards[idx]
	}

	return nil
}
