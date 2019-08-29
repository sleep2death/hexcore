package cards

import (
	"errors"
	"math/rand"
	"strconv"
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

	SetName(name string)
	Name() string

	SetID(id string)
	ID() string

	Copy() Card
}

// Pile of the cards
type Pile []Card

// Shuffle the cards
func (p *Pile) Shuffle() {
	rand.Shuffle(len(*p), func(i, j int) { (*p)[i], (*p)[j] = (*p)[j], (*p)[i] })
}

// Draw one card from the source pile
func (p *Pile) Draw(source *Pile) (Card, error) {
	if len(*source) == 0 {
		return nil, ErrNotEnoughCards
	}

	*p = append(*p, (*source)[len(*source)-1])
	*source = (*source)[:len(*source)-1]

	card := (*p)[len(*p)-1]
	return card, nil
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
func (p *Pile) Pick(id string, source *Pile) (Card, error) {
	card, idx, err := source.FindCard(id)
	if err != nil {
		return nil, err
	}

	*p = append(*p, card)

	// card = (*source)[idx]
	copy((*source)[idx:], (*source)[idx+1:])
	*source = (*source)[:len(*source)-1]

	return card, nil
}

// Copy every card of the source pile
func (p *Pile) Copy() *Pile {
	copy := make(Pile, 0)
	for _, card := range *p {
		copy = append(copy, card.Copy())
	}
	return &copy
}

// TestCard -
type TestCard struct {
	id     string
	name   string
	copied int
}

func (c *TestCard) String() string {
	return "<card " + c.id + ">"
}

// SetName -
func (c *TestCard) SetName(name string) {
	c.name = name
}

// Name -
func (c *TestCard) Name() string {
	return c.name
}

// SetID -
func (c *TestCard) SetID(id string) {
	c.id = id
}

// ID -
func (c *TestCard) ID() string {
	return c.id
}

// Upgrade -
func (c *TestCard) Upgrade() error {
	return errors.New("can't upgrade")
}

// Copy -
func (c *TestCard) Copy() Card {
	c.copied++
	return &TestCard{
		id: "copy:" + strconv.Itoa(c.copied) + " of <" + c.id + ">",
	}
}
