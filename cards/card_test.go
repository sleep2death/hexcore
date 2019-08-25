package cards

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCard(t *testing.T) {
	cardA := CreateCardFunc["Strike"]()
	assert.EqualValues(t, cardA.Name(), "Strike")

	// Upgrade the card
	cardA.UpgradeBase()
	assert.Equal(t, cardA.Base().Level, 1)      // level from 0 to 1
	assert.Equal(t, cardA.Base().Damage, 9)     // damage from 6 to 9
	assert.Equal(t, Enemy, cardA.Base().Target) // if upgrade target not set, then not change it

	assert.Equal(t, cardA.Current().Level, 1) // upgrade base will also update the current

	cardA.UpgradeCurrent()
	assert.Equal(t, cardA.Current().Level, 2) // upgrade the current
	assert.Equal(t, cardA.Base().Level, 1)    // upgrade current will NOT update the base

	// Copy the Card
	cardACopy := cardA.Copy()
	assert.Equal(t, 1, cardACopy.Base().Level)    // base status copied from cardA
	assert.Equal(t, 1, cardACopy.Current().Level) // current status will NOT be copied from cardA, it's a copy of the base

	cardB := CreateCardFunc["Defend"]()
	assert.Equal(t, cardB.Name(), "Defend")
	assert.Equal(t, cardB.Base().Damage, 0)
	assert.Equal(t, cardB.Base().Block, 5)

	cardB.UpgradeCurrent()
	assert.Equal(t, cardB.Current().Level, 1)     // level from 0 to 1
	assert.Equal(t, cardB.Current().Block, 8)     // damage from 6 to 9
	assert.Equal(t, cardB.Current().Target, Self) // if upgrade target not set, then not change it
}

func TestPile(t *testing.T) {
	random := rand.New(rand.NewSource(9012))
	p, err := CreatePile([]string{"Strike", "Strike", "Strike", "Strike", "Strike"})

	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, 5, p.Num())

	// Draw
	pp, err := CreatePile([]string{"Defend", "Defend", "Defend", "Defend", "Bash"})
	pp.Draw(p, p.Num())
	assert.Equal(t, 10, len(pp.cards))
	assert.Equal(t, 0, len(p.cards))
	assert.EqualError(t, pp.Draw(p, 1), "not enough card(s) to draw")
	assert.EqualError(t, pp.Draw(p, 0), "draw number should be larger than 0")

	// Pick
	p.Pick(pp, pp.cards[0].ID())
	assert.Equal(t, 1, p.Num())
	assert.Equal(t, 9, pp.Num())
	assert.Equal(t, "Defend", p.cards[0].Name())
	p.Pick(pp, pp.cards[0].ID())
	assert.Equal(t, 2, p.Num())
	assert.Equal(t, 8, pp.Num())
	assert.Equal(t, "Defend", p.cards[1].Name())

	// Shuffle
	pp.Shuffle(random)
	assert.Equal(t, "Strike", pp.cards[0].Name())
	assert.Equal(t, "Defend", pp.cards[7].Name())
}

// 	pp, err := CreatePile(random, []string{"Defend", "Defend", "Defend", "Defend", "Bash"})

// 	pp.Draw(len(pp.cards), p)

// 	assert.Equal(t, 10, len(p.cards))

// 	assert.Equal(t, 0, len(pp.cards))

// 	assert.EqualError(t, pp.Draw(1, p), "not enough card(s) to draw")

// 	assert.EqualError(t, pp.Draw(0, p), "draw number should be larger than 0")

// init all the cards in the pile

// 	for _, card := range p.cards {
// 		card.Init()
// 	}
// 	assert.Equal(t, 2, p.FindCardByID(p.cards[2].ID()))
// 	assert.Equal(t, 10, p.CardsNum())

// 	card, err := p.RemoveCard(2)
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	assert.Equal(t, 9, p.CardsNum())
// 	assert.Equal(t, "[strike]", card.String())

// 	card, err = p.RemoveCard(2)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	assert.Equal(t, 8, p.CardsNum())
// 	assert.Equal(t, "[defend]", card.String())

// 	card, err = p.RemoveCard(7)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	assert.Equal(t, 7, p.CardsNum())
// 	assert.Equal(t, "[bash]", card.String())
// }

// func TestStatusCopy(t *testing.T) {

// 	c := CreateCardFunc["Strike"]()

// 	s := c.Base().Copy()

// 	assert.Equal(t, c.Base(), s)

// 	t.Logf("origin -> %v, copied -> %v", c.Base(), s)

// }

// func TestCardCopy(t *testing.T) {
// 	c := CreateCardFunc["Strike"]()
// 	c.Init()
// 	s := c.Copy()

// 	// the base status pointer should be the same
// 	assert.Equal(t, fmt.Sprintf("%p", c.Base()), fmt.Sprintf("%p", s.Base()))

// 	// the current status pointer should not be the same
// 	assert.NotEqual(t, fmt.Sprintf("%p", c.Current()), fmt.Sprintf("%p", s.Current()))

// 	// like card [Ritual Dagger] -  if this card kills an enemy then permanently increase this card's damage by 3(5)
// 	// if card[Ritual Dagger] upgraded in the battle, then original card in the deck will also be upgraded
// 	// manager can use "base" status permanently change the card
// 	s.Base().Damage = 100
// 	assert.Equal(t, 100, c.Base().Damage)
// }

// // func TestManager(t *testing.T) {

// // 	m, err := CreateManager([]string{"Strike", "Strike", "Strike", "Strike", "Strike", "Defend", "Defend", "Defend", "Defend", "Bash"})

// // 	if err != nil {
// // 		t.Error(err)
// // 	}

// // 	assert.Equal(t, 10, len(m.deck.cards))

// // 	// if a card name not exist, then return an error, and clear the pile

// // 	m, err = CreateManager([]string{"Strike", "Strike", "Strike", "Strike", "Strike", "Defend", "Defend", "ABC", "Defend", "Bash"})

// // 	assert.EqualError(t, err, "create function for card [ABC] not found")

// // 	// assert.Equal(t, 0, len(m.deck.cards))
// // 	// create the deck

// // 	m, err = CreateManager([]string{"Strike", "Strike", "Strike", "Strike", "Strike", "Defend", "Defend", "Defend", "Defend", "Bash"})

// // 	if err != nil {
// // 		t.Error(err)
// // 	}

// // 	// copy the deck cards into draw pile, and shuffle
// // 	err = m.Shuffle()

// // 	if err != nil {
// // 		t.Error(err)
// // 	}

// // 	// t.Log(m.deck.cards, m.draw.cards)

// // 	// [[bash] [defend] [defend] [strike] [strike] [strike] [strike] [defend] [defend] [strike]
// // 	// draw 5 cards from draw pile to hand pile
// // 	err = m.Draw(5)
// // 	if err != nil {
// // 		t.Error(err)
// // 	}

// // 	assert.Equal(t, 5, m.draw.CardsNum()) // [[bash] [defend] [defend] [strike] [strike]]
// // 	assert.Equal(t, 5, m.hand.CardsNum()) // [[strike] [strike] [defend] [defend] [strike]]

// // 	// exaust a card from hand pile
// // 	eCard := m.hand.cards[2]
// // 	err = m.Exaust(eCard)
// // 	if err != nil {
// // 		t.Error(err)
// // 	}
// // 	assert.Equal(t, 1, m.exaust.CardsNum()) // [[defend]]
// // 	assert.Equal(t, 4, m.hand.CardsNum())   // [[strike] [strike] [defend] [strike]]
// // 	assert.Equal(t, "[defend]", eCard.String())
// // 	assert.Equal(t, m.exaust.cards[0], eCard)

// // 	// discard a card from hand pile
// // 	dCard := m.hand.cards[2]
// // 	err = m.Discard(dCard)
// // 	if err != nil {
// // 		t.Error(err)
// // 	}
// // 	assert.Equal(t, 1, m.discard.CardsNum()) // [[defend]]
// // 	assert.Equal(t, 3, m.hand.CardsNum())    // [[strike] [strike] [strike]]
// // 	assert.Equal(t, "[defend]", dCard.String())
// // 	assert.Equal(t, m.discard.cards[0], dCard)

// // 	// create the deck
// // 	m, err = CreateManager([]string{"Strike", "Strike", "Strike", "Strike", "Strike", "Defend", "Defend", "Defend", "Defend", "Bash"})
// // 	if err != nil {
// // 		t.Error(err)
// // 	}

// // 	// copy the deck cards into draw pile, and shuffle
// // 	err = m.Shuffle()
// // 	if err != nil {
// // 		t.Error(err)
// // 	}

// // 	assert.Equal(t, 10, m.draw.CardsNum()) // [[bash] [defend] [defend] [strike] [strike] [strike] [strike] [defend] [defend] [strike]]

// // 	err = m.Draw(5)
// // 	if err != nil {
// // 		t.Error(err)
// // 	}

// // 	assert.Equal(t, 5, m.hand.CardsNum()) // hand: [strike] [strike] [defend] [defend] [strike]
// // 	err = m.ExaustAll()
// // 	if err != nil {
// // 		t.Error(err)
// // 	}
// // 	assert.Equal(t, 0, m.hand.CardsNum())
// // 	assert.Equal(t, 5, m.exaust.CardsNum())

// // 	err = m.Draw(5)
// // 	assert.Equal(t, 5, m.hand.CardsNum()) // hand: [bash] [defend] [defend] [strike] [strike]
// // 	err = m.DiscardAll()
// // 	if err != nil {
// // 		t.Error(err)
// // 	}
// // 	assert.Equal(t, 0, m.draw.CardsNum())
// // 	assert.Equal(t, 0, m.hand.CardsNum())

// // 	assert.Equal(t, 5, m.exaust.CardsNum())
// // 	assert.Equal(t, 5, m.discard.CardsNum())

// // 	if err = m.ReShuffle(); err != nil {
// // 		t.Error(err)
// // 	}

// // 	assert.Equal(t, 0, m.discard.CardsNum())
// // 	assert.Equal(t, 5, m.draw.CardsNum())
// // }
