package hexcore

import (
	"testing"

	"github.com/sleep2death/hexcore/cards"
)

func TestBattleInit(t *testing.T) {
	cards, err := cards.CreateManager([]string{"Strike", "Strike", "Strike", "Strike", "Strike", "Defend", "Defend", "Defend", "Defend", "Bash"})
	if err != nil {
		t.Error(err)
	}
	b := CreateBattle(cards)

	if err := b.cards.Shuffle(); err != nil {
		t.Error(err)
	}

	// deck: [[strike] [strike] [strike] [strike] [strike] [defend] [defend] [defend] [defend] [bash]]
	// draw: [[bash] [defend] [defend] [strike] [strike] [strike] [strike] [defend] [defend] [strike]]

	if err := b.cards.Draw(5); err != nil {
		t.Error(err)
	}

	b.WaitForTarget()

	// deck: [[strike] [strike] [strike] [strike] [strike] [defend] [defend] [defend] [defend] [bash]]
	// draw: [[bash] [defend] [defend] [strike] [strike]]
	// hand: [[strike] [strike] [defend] [defend] [strike]]
}
