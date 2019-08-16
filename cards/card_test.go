package cards

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCardCreate(t *testing.T) {
	card := CreateCardFunc["Strike"]()
	assert.EqualValues(t, card.String(), "[strike]")
	assert.EqualValues(t, card.Base().Actions.Play[0], DealDamage)
	// Upgrade the card
	card.Upgrade()
	assert.EqualValues(t, card.Base().Level, 1)           // level from 0 to 1
	assert.EqualValues(t, card.Base().Damage, 9)          // damage from 6 to 9
	assert.Equal(t, "enemy", card.Base().Target.String()) // if upgrade target not set, then not change it
}

func TestPile(t *testing.T) {
	p := &Pile{
		cards: make([]Card, 0, 50),
	}

	// create a defaut pile
	p.AddToTop(
		CreateCardFunc["Strike"](),
		CreateCardFunc["Strike"](),
		CreateCardFunc["Strike"](),
		CreateCardFunc["Strike"](),
		CreateCardFunc["Strike"](),
	)

	pp := &Pile{
		cards: make([]Card, 0, 50),
	}

	pp.AddToTop(
		CreateCardFunc["Defend"](),
		CreateCardFunc["Defend"](),
		CreateCardFunc["Defend"](),
		CreateCardFunc["Defend"](),
		CreateCardFunc["Bash"](),
	)

	pp.Draw(len(pp.cards), p)

	assert.Equal(t, 10, len(p.cards))
	assert.Equal(t, 0, len(pp.cards))

	assert.EqualError(t, pp.Draw(1, p), "not enough card(s) to draw")
	assert.EqualError(t, pp.Draw(0, p), "n(0) should be larger than 0")

	t.Logf("%v", p.cards)

	// shuffle the pile with a given static number seed
	p.Shuffle(5)
	t.Logf("%v", p.cards)
	assert.Equal(t, p.cards[0].String(), "[defend]")
}
