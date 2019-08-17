package cards

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCardCreate(t *testing.T) {

	card := CreateCardFunc["Strike"]()

	assert.EqualValues(t, card.String(), "[strike]")

	assert.EqualValues(t, card.Base().Actions.Play[0], DealDamage)

	// Upgrade the card

	card.Upgrade()

	assert.EqualValues(t, card.Base().Level, 1) // level from 0 to 1

	assert.EqualValues(t, card.Base().Damage, 9) // damage from 6 to 9

	assert.Equal(t, "enemy", card.Base().Target.String()) // if upgrade target not set, then not change it

}

func TestPile(t *testing.T) {
	p := &Pile{}
	p.CreateCardByName([]string{"Strike", "Strike", "Strike", "Strike", "Strike"})

	pp := &Pile{}
	pp.CreateCardByName([]string{"Defend", "Defend", "Defend", "Defend", "Bash"})

	pp.Draw(len(pp.cards), p)

	assert.Equal(t, 10, len(p.cards))

	assert.Equal(t, 0, len(pp.cards))

	assert.EqualError(t, pp.Draw(1, p), "not enough card(s) to draw")

	assert.EqualError(t, pp.Draw(0, p), "n(0) should be larger than 0")

	t.Logf("%v", p.cards)

	// shuffle the pile with a given static number seed

	p.Shuffle()

	t.Logf("%v", p.cards)

	assert.Equal(t, p.cards[0].String(), "[strike]")

	// init all the cards in the pile

	for _, card := range p.cards {

		card.Init()

	}

	assert.Equal(t, 2, p.FindCardByID(p.cards[2].ID()))

}

func TestStatusCopy(t *testing.T) {

	c := CreateCardFunc["Strike"]()

	s := c.Base().Copy()

	assert.Equal(t, c.Base(), s)

	t.Logf("origin -> %v, copied -> %v", c.Base(), s)

}

func TestCardCopy(t *testing.T) {

	c := CreateCardFunc["Strike"]()

	c.Init()

	s := c.Copy()

	// the base status pointer should be the same

	assert.Equal(t, fmt.Sprintf("%p", c.Base()), fmt.Sprintf("%p", s.Base()))

	// the current status pointer should not be the same

	assert.NotEqual(t, fmt.Sprintf("%p", c.Current()), fmt.Sprintf("%p", s.Current()))

	assert.NotEqual(t, fmt.Sprintf("%p", c.Current().Actions), fmt.Sprintf("%p", s.Current().Actions))

	// like card [Ritual Dagger] -  if this card kills an enemy then permanently increase this card's damage by 3(5)

	// if card[Ritual Dagger] upgraded in the battle, then original card in the deck will also be upgraded

	// manager can use "base" status permanently change the card

	s.Base().Damage = 100

	assert.Equal(t, 100, c.Base().Damage)

}

func TestManager(t *testing.T) {

	m := &Manager{}

	err := m.Init([]string{"Strike", "Strike", "Strike", "Strike", "Strike", "Defend", "Defend", "Defend", "Defend", "Bash"})

	if err != nil {

		t.Error(err)

	}

	assert.Equal(t, 10, len(m.deck.cards))

	// if a card name not exist, then return an error, and clear the pile

	err = m.Init([]string{"Strike", "Strike", "Strike", "Strike", "Strike", "Defend", "Defend", "ABC", "Defend", "Bash"})

	assert.EqualError(t, err, "create function for card [ABC] not found")

	assert.Equal(t, 0, len(m.deck.cards))

	// create the deck

	err = m.Init([]string{"Strike", "Strike", "Strike", "Strike", "Strike", "Defend", "Defend", "Defend", "Defend", "Bash"})

	if err != nil {

		t.Error(err)

	}

	// copy the deck cards into draw pile, and shuffle

	err = m.Shuffle()

	if err != nil {

		t.Error(err)

	}

	// t.Log(m.deck.cards, m.draw.cards)

	// [[bash] [defend] [defend] [strike] [strike] [strike] [strike] [defend] [defend] [strike]

	err = m.Draw(5)

	if err != nil {

		t.Error(err)

	}

	assert.Equal(t, 5, m.draw.CardsNum()) // [[bash] [defend] [defend] [strike] [strike]]

	assert.Equal(t, 5, m.hand.CardsNum()) // [[strike] [strike] [defend] [defend] [strike]]

	t.Log(m.draw.cards)

	t.Log(m.hand.cards)

}