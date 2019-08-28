package cards

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShuffle(t *testing.T) {
	cards := Pile([]Card{
		&TestCard{id: "a"},
		&TestCard{id: "b"},
		&TestCard{id: "c"},
		&TestCard{id: "d"},
		&TestCard{id: "e"},
		&TestCard{id: "f"},
	})

	cards[3].SetID("3")
	seed := rand.New(rand.NewSource(99))
	cards.Shuffle(seed)
	assert.Equal(t, "[<card f> <card a> <card b> <card c> <card e> <card 3>]", fmt.Sprint(cards))
}

func TestDraw(t *testing.T) {
	a := Pile([]Card{
		&TestCard{id: "a"},
		&TestCard{id: "b"},
		&TestCard{id: "c"},
	})

	b := Pile([]Card{
		&TestCard{id: "d"},
		&TestCard{id: "e"},
		&TestCard{id: "f"},
	})

	a.Draw(&b)
	assert.Equal(t, "[<card a> <card b> <card c> <card f>]", fmt.Sprint(a))
	assert.Equal(t, "[<card d> <card e>]", fmt.Sprint(b))

	a.Draw(&b)
	a.Draw(&b)
	assert.Equal(t, "[<card a> <card b> <card c> <card f> <card e> <card d>]", fmt.Sprint(a))

	_, err := a.Draw(&b)
	assert.Equal(t, ErrNotEnoughCards, err)

}

func TestPick(t *testing.T) {
	a := Pile([]Card{
		&TestCard{id: "a"},
		&TestCard{id: "b"},
		&TestCard{id: "c"},
		&TestCard{id: "d"},
		&TestCard{id: "e"},
		&TestCard{id: "f"},
	})

	b := Pile([]Card{
		&TestCard{id: "g"},
		&TestCard{id: "h"},
	})

	c := Pile([]Card{})

	card, idx, _ := a.FindCard("a")
	assert.Equal(t, "<card a>", card.String())
	assert.Equal(t, 0, idx)

	b.Pick("d", &a)
	assert.Equal(t, "[<card g> <card h> <card d>]", fmt.Sprint(b))
	assert.Equal(t, "[<card a> <card b> <card c> <card e> <card f>]", fmt.Sprint(a))

	_, err := b.Pick("d", &a)
	assert.Equal(t, ErrCardNotExist, err)

	_, err = b.Pick("d", &c)
	assert.Equal(t, ErrPileIsNilOrEmpty, err)
}

func TestCopy(t *testing.T) {
	a := &TestCard{id: "a"}
	b := a.Copy()
	c := a.Copy()
	d := b.Copy()
	assert.Equal(t, "copy:1 of <a>", b.ID())
	assert.Equal(t, "copy:2 of <a>", c.ID())
	assert.Equal(t, "copy:1 of <copy:1 of <a>>", d.ID())

	err := a.Upgrade()
	assert.Equal(t, "can't upgrade", err.Error())

	p := Pile([]Card{
		&TestCard{id: "a"},
		&TestCard{id: "b"},
		&TestCard{id: "c"},
	})

	copy := p.Copy()
	assert.Equal(t, 3, len(*copy))
	assert.Equal(t, "&[<card copy:1 of <a>> <card copy:1 of <b>> <card copy:1 of <c>>]", fmt.Sprint(copy))
}
