package cards

import (
	"fmt"
	"math/rand"
)

// Pile of cards
type Pile struct {
	seed  int64
	cards []Card
}

// AddToTop with the given card(s)
func (p *Pile) AddToTop(c ...Card) {
	p.cards = append(p.cards, c...)
}

// AddToBottom with the given card(s)
func (p *Pile) AddToBottom(c ...Card) {
	p.cards = append(c, p.cards...)
}

// Draw n card(s) to the target pile
func (p *Pile) Draw(n int, target *Pile) error {
	if n <= 0 {
		return errDrawNumber
	}
	if n > len(p.cards) {
		return errNotEnoughCards
	}

	idx := len(p.cards) - n

	target.cards = append(target.cards, p.cards[idx:]...)
	p.cards = p.cards[:idx]
	return nil
}

// DrawCard draw one card by the given index to the target pile
func (p *Pile) DrawCard(i int, target *Pile) error {
	if i < 0 || i > len(p.cards)-1 {
		return errDrawIndex
	}

	card := p.cards[i]
	p.RemoveCard(i)

	target.AddToTop(card)

	return nil
}

// RemoveCard from the pile
func (p *Pile) RemoveCard(i int) error {
	if i < 0 || i > len(p.cards)-1 {
		return errDrawIndex
	}
	copy(p.cards[i:], p.cards[i+1:])
	p.cards = p.cards[:len(p.cards)-1]
	return nil
}

// FindCardByID return the card index with given id
func (p *Pile) FindCardByID(id string) int {
	if p.CardsNum() == 0 {
		return -1
	}

	for i, c := range p.cards {
		if c.ID() == id {
			return i
		}
	}

	return -1
}

// Shuffle the pile
func (p *Pile) Shuffle() {
	if p.CardsNum() <= 0 {
		return
	}

	rand.Seed(p.seed)
	rand.Shuffle(len(p.cards), func(i, j int) { p.cards[i], p.cards[j] = p.cards[j], p.cards[i] })
}

// CardsNum - get the card number of the pile
func (p *Pile) CardsNum() int {
	return len(p.cards)
}

// CreateCardByName - create the card by the given name
func (p *Pile) CreateCardByName(cardSet []string) error {
	for _, s := range cardSet {
		if CreateCardFunc[s] == nil {
			// clear all the items by setting the slice to nil
			// see: https://stackoverflow.com/questions/16971741/how-do-you-clear-a-slice-in-go
			p.cards = nil
			return fmt.Errorf("create function for card [%s] not found", s)
		}

		card := CreateCardFunc[s]()
		p.AddToTop(card)
	}
	return nil
}
