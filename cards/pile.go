package cards

import (
	"errors"
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
		return fmt.Errorf("n(%d) should be larger than 0", n)
	}
	if n > len(p.cards) {
		return errors.New("not enough card(s) to draw")
	}

	idx := len(p.cards) - n

	target.cards = append(target.cards, p.cards[idx:]...)
	p.cards = p.cards[:idx]
	return nil
}

// FindCardByID return the card index with given id
func (p *Pile) FindCardByID(id string) int {
	for i, c := range p.cards {
		if c.ID() == id {
			return i
		}
	}

	return -1
}

// Shuffle the pile
func (p *Pile) Shuffle() {
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
