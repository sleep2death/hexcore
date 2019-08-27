package cards

import (
	"errors"
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestCard struct {
	id   string
	desc string
}

func (c *TestCard) Upgrade() error {
	return errors.New("can't upgrade")
}

func (c *TestCard) String() string {
	return "<" + c.desc + ">"
}

func (c *TestCard) ID() string {
	return c.id
}

func TestShuffle(t *testing.T) {
	cards := Pile([]Card{
		&TestCard{id: "a", desc: "card a"},
		&TestCard{id: "b", desc: "card b"},
		&TestCard{id: "c", desc: "card c"},
		&TestCard{id: "d", desc: "card d"},
		&TestCard{id: "e", desc: "card e"},
		&TestCard{id: "f", desc: "card f"},
	})
	seed := rand.New(rand.NewSource(99))
	cards.Shuffle(seed)
	assert.Equal(t, "[<card f> <card a> <card b> <card c> <card e> <card d>]", fmt.Sprint(cards))
}

func TestDraw(t *testing.T) {
	a := Pile([]Card{
		&TestCard{id: "a", desc: "card a"},
		&TestCard{id: "b", desc: "card b"},
		&TestCard{id: "c", desc: "card c"},
	})

	b := Pile([]Card{
		&TestCard{id: "d", desc: "card d"},
		&TestCard{id: "e", desc: "card e"},
		&TestCard{id: "f", desc: "card f"},
	})

	a.Draw(&b)
	assert.Equal(t, "[<card a> <card b> <card c> <card f>]", fmt.Sprint(a))
	assert.Equal(t, "[<card d> <card e>]", fmt.Sprint(b))

	a.Draw(&b)
	a.Draw(&b)
	assert.Equal(t, "[<card a> <card b> <card c> <card f> <card e> <card d>]", fmt.Sprint(a))

	err := a.Draw(&b)
	assert.Equal(t, ErrNotEnoughCards, err)
}

func TestPick(t *testing.T) {
	a := Pile([]Card{
		&TestCard{id: "a", desc: "card a"},
		&TestCard{id: "b", desc: "card b"},
		&TestCard{id: "c", desc: "card c"},
		&TestCard{id: "d", desc: "card d"},
		&TestCard{id: "e", desc: "card e"},
		&TestCard{id: "f", desc: "card f"},
	})

	b := Pile([]Card{
		&TestCard{id: "g", desc: "card g"},
		&TestCard{id: "h", desc: "card h"},
	})

	c := Pile([]Card{})

	card, idx, _ := a.FindCard("a")
	assert.Equal(t, "<card a>", card.String())
	assert.Equal(t, 0, idx)

	b.Pick(&a, "d")
	assert.Equal(t, "[<card g> <card h> <card d>]", fmt.Sprint(b))
	assert.Equal(t, "[<card a> <card b> <card c> <card e> <card f>]", fmt.Sprint(a))

	err := b.Pick(&a, "d")
	assert.Equal(t, ErrCardNotExist, err)

	err = b.Pick(&c, "d")
	assert.Equal(t, ErrPileIsNilOrEmpty, err)
}
