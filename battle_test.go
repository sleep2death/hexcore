package hexcore

import (
	"testing"

	"github.com/sleep2death/hexcore/cards"
	"gopkg.in/go-playground/assert.v1"
)

func TestBattleInit(t *testing.T) {
	b := &Battle{cards: &cards.Manager{}}
	assert.Equal(t, cards.ErrPileIsNil, b.cards.Shuffle())

	if err := b.cards.Create([]string{"Strike", "Strike", "Strike", "Strike", "Strike", "Defend", "Defend", "Defend", "Defend", "Bash"}); err != nil {
		t.Error(err)
	}

	if err := b.cards.Shuffle(); err != nil {
		t.Error(err)
	}

	// deck: [[strike] [strike] [strike] [strike] [strike] [defend] [defend] [defend] [defend] [bash]]
	// draw: [[bash] [defend] [defend] [strike] [strike] [strike] [strike] [defend] [defend] [strike]]

	if err := b.cards.Draw(5); err != nil {
		t.Error(err)
	}

	// deck: [[strike] [strike] [strike] [strike] [strike] [defend] [defend] [defend] [defend] [bash]]
	// draw: [[bash] [defend] [defend] [strike] [strike]]
	// hand: [[strike] [strike] [defend] [defend] [strike]]
}
