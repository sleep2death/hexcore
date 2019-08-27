package cards

import (
	"errors"
	"math/rand"
)

var (
	// ErrDrawNumber -
	ErrDrawNumber = errors.New("draw number should be larger than 0")

	// ErrDrawIndex -
	ErrDrawIndex = errors.New("draw index should be larger than 0 and less than len(cards) - 1")

	// ErrNotEnoughCards -
	ErrNotEnoughCards = errors.New("not enough card(s) to draw")

	// ErrCardNotExist -
	ErrCardNotExist = errors.New("card doesn't exist")

	// ErrPileIsNilOrEmpty -
	ErrPileIsNilOrEmpty = errors.New("target pile is nil or empty")
)

// Card - interface
type Card interface {
	Upgrade() error
	String() string
	ID() string
}

// Pile of the cards
type Pile []Card

// Shuffle the cards
func (p *Pile) Shuffle(seed *rand.Rand) {
	rand.Shuffle(len(*p), func(i, j int) { (*p)[i], (*p)[j] = (*p)[j], (*p)[i] })
}

// Draw one card from the source pile
func (p *Pile) Draw(source *Pile) error {
	if len(*source) == 0 {
		return ErrNotEnoughCards
	}

	*p = append(*p, (*source)[len(*source)-1])
	*source = (*source)[:len(*source)-1]
	return nil
}

// FindCard return both the card and card index of the pile
func (p *Pile) FindCard(id string) (card Card, idx int, err error) {
	if len(*p) == 0 {
		return nil, -1, ErrPileIsNilOrEmpty
	}

	for i, c := range *p {
		if c.ID() == id {
			return c, i, nil
		}
	}

	return nil, -1, ErrCardNotExist
}

// Pick one card from source pile and add it to the top of the pile
func (p *Pile) Pick(source *Pile, id string) error {
	card, idx, err := source.FindCard(id)
	if err != nil {
		return err
	}

	*p = append(*p, card)

	card = (*source)[idx]
	copy((*source)[idx:], (*source)[idx+1:])
	*source = (*source)[:len(*source)-1]

	return nil
}
