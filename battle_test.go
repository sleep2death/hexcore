package hexcore

import (
	"testing"

	"github.com/sleep2death/hexcore/cards"
	"github.com/stretchr/testify/assert"
)

func TestBattleInit(t *testing.T) {
	manager, err := cards.CreateManager([]string{"Strike", "Strike", "Strike", "Strike", "Strike", "Defend", "Defend", "Defend", "Defend", "Bash"})

	if err != nil {
		t.Error(err)
	}
	b := CreateBattle(manager)

	if err := b.cards.Shuffle(); err != nil {
		t.Error(err)
	}

	// deck: [[strike] [strike] [strike] [strike] [strike] [defend] [defend] [defend] [defend] [bash]]
	// draw: [[bash] [defend] [defend] [strike] [strike] [strike] [strike] [defend] [defend] [strike]]

	if err := b.cards.Draw(5); err != nil {
		t.Error(err)
	}

	go b.Update()

	// t.Log(b.cards)

	assert.Equal(t, b.InputState(), WaitForNone)

	b.SetInputState(WaitForCardPlay)

	b.CardInput <- b.cards.GetCardByIndex(0, cards.Hand).ID() // First Card - "Strike" in the hand pile
	b.CardInput <- b.cards.GetCardByIndex(1, cards.Hand).ID() // First Card - "Strike" in the hand pile
	b.CardInput <- b.cards.GetCardByIndex(2, cards.Hand).ID() // First Card - "Strike" in the hand pile
	// b.CardInput <- "f24XuufQVDpQdf5qmLNZjR"                   // First Card - "Strike" in the hand pile
	// b.CardInput <- "f24XuufQVDpQdf5qmLNZjR"                   // First Card - "Strike" in the hand pile

	// t.Log(b.cards)

	// deck: [[strike] [strike] [strike] [strike] [strike] [defend] [defend] [defend] [defend] [bash]]
	// draw: [[bash] [defend] [defend] [strike] [strike]]
	// hand: [[strike] [strike] [defend] [defend] [strike]]
}
