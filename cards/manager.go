package cards

import (
	"errors"
)

var seed int64 = 9012

var errDrawNumber error = errors.New("draw number should be larger than 0")
var errDrawIndex error = errors.New("draw index should be larger than 0 and less than len(cards) - 1")
var errNotEnoughCards error = errors.New("not enough card(s) to draw")
var errCardNotExist error = errors.New("card doesn't exist in the pile")

// Manager of all the cards in player's pocket
type Manager struct {
	deck    *Pile
	draw    *Pile
	hand    *Pile
	discard *Pile
	exaust  *Pile
}

// Init the card manager by filling the deck with given cards
func (m *Manager) Init(cardSet []string) error {
	cap := len(cardSet)
	dCap := cap * 4
	m.deck = &Pile{seed: seed, cards: make([]Card, 0, cap)}
	m.draw = &Pile{seed: seed, cards: make([]Card, 0, dCap)}
	m.hand = &Pile{seed: seed, cards: make([]Card, 0, dCap)}
	m.discard = &Pile{seed: seed, cards: make([]Card, 0, dCap)}
	m.exaust = &Pile{seed: seed, cards: make([]Card, 0, dCap)}
	return m.deck.CreateCardByName(cardSet)
}

// Shuffle - copy all cards from the deck pile to draw pile
func (m *Manager) Shuffle() error {
	if m.deck == nil || m.deck.CardsNum() <= 0 {
		return errors.New("empty deck pile")
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
		return errCardNotExist
	}
	return nil
}

// Discard cards from hand pile to discard pile
func (m *Manager) Discard(card Card) {
}
