package hexcore

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/sleep2death/hexcore/actions"

	"github.com/sleep2death/hexcore/cards"
)

func TestBattleInit(t *testing.T) {
	seed := rand.New(rand.NewSource(9012))
	deck, _ := cards.CreatePile(seed, []string{"Strike", "Strike", "Strike", "Strike", "Strike", "Defend", "Defend", "Defend", "Defend", "Bash"})

	battle := CreateBattle(deck)
	go battle.Update()

	var action actions.Action

	action = &actions.StartBattleAction{
		Deck: battle.cards[DeckPile],
		Draw: battle.cards[DrawPile],
	}

	go battle.Exec([]actions.Action{action})
	time.Sleep(time.Millisecond * 10)

	assert.Equal(t, 10, battle.cards[DrawPile].Num())
	assert.Equal(t, "[Bash Defend Defend Strike Strike Strike Strike Defend Defend Strike]", battle.cards[DrawPile].String())

	action = &actions.StartTurnAction{
		Draw:    battle.cards[DrawPile],
		Hand:    battle.cards[HandPile],
		Discard: battle.cards[DiscardPile],
	}

	go battle.Exec([]actions.Action{action})
	time.Sleep(time.Millisecond * 10)

	assert.Equal(t, 5, battle.cards[DrawPile].Num())
	assert.Equal(t, 5, battle.cards[HandPile].Num())
	assert.Equal(t, "[Bash Defend Defend Strike Strike]", battle.cards[DrawPile].String())
	assert.Equal(t, "[Strike Strike Defend Defend Strike]", battle.cards[HandPile].String())

	// Discard 5 times, then hand pile shoud be empty
	discardBatch := make([]actions.Action, 5)
	for i := 0; i < 5; i++ {
		card, _ := battle.cards[HandPile].GetCard(i)
		discardBatch[i] = &actions.DiscardAction{
			Hand:    battle.cards[HandPile],
			Discard: battle.cards[DiscardPile],
			ID:      card.ID(),
		}
	}

	go battle.Exec(discardBatch)
	time.Sleep(time.Millisecond * 10)

	assert.Equal(t, 5, battle.cards[DrawPile].Num())
	assert.Equal(t, 0, battle.cards[HandPile].Num())
	assert.Equal(t, "[Bash Defend Defend Strike Strike]", battle.cards[DrawPile].String())
	assert.Equal(t, "[]", battle.cards[HandPile].String())
	assert.Equal(t, "[Strike Strike Defend Defend Strike]", battle.cards[DiscardPile].String())

	// Discard an empty hand pile should return error
	action = &actions.DiscardAllAction{
		Hand:    battle.cards[HandPile],
		Discard: battle.cards[DiscardPile],
	}

	go battle.Exec([]actions.Action{action})
	time.Sleep(time.Millisecond * 10)

	// draw 3 more cards from draw pile into hand pile
	// and discard them again
	go battle.Exec([]actions.Action{
		&actions.DrawAction{
			Pile:     battle.cards[HandPile],
			DrawFrom: battle.cards[DrawPile],
			N:        3,
		},
		&actions.DiscardAllAction{
			Hand:    battle.cards[HandPile],
			Discard: battle.cards[DiscardPile],
		},
	})
	time.Sleep(time.Millisecond * 10)

	assert.Equal(t, "[Bash Defend]", battle.cards[DrawPile].String())
	assert.Equal(t, "[]", battle.cards[HandPile].String())
	assert.Equal(t, "[Strike Strike Defend Defend Strike Defend Strike Strike]", battle.cards[DiscardPile].String())

	// Start turn again, draw pile should be empty
	action = &actions.StartTurnAction{
		Draw:    battle.cards[DrawPile],
		Hand:    battle.cards[HandPile],
		Discard: battle.cards[DiscardPile],
	}

	go battle.Exec([]actions.Action{action})
	time.Sleep(time.Millisecond * 10)

	action = &actions.WaitForPlayAction{
		Hand:  battle.cards[HandPile],
		Input: battle.input,
	}

	go battle.Exec([]actions.Action{action})
	time.Sleep(time.Second * 6)
}
