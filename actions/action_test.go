package actions

import (
	"math/rand"
	"testing"
	"time"

	"github.com/sleep2death/hexcore/cards"
	"github.com/stretchr/testify/assert"
)

func TestStartBattleAction(t *testing.T) {
	seed := rand.New(rand.NewSource(9012))
	errorc := make(chan error)
	output := make(chan Action)

	deck, _ := cards.CreatePile(seed, []string{"Strike", "Strike", "Strike", "Strike", "Strike", "Defend", "Defend", "Defend", "Defend", "Bash"})
	draw, _ := cards.CreatePile(seed, nil)

	action := &StartBattleAction{
		Deck: deck,
		Draw: draw,
	}

	go Execute([]Action{action}, output, errorc)
	time.Sleep(time.Millisecond * 10)
	assert.Equal(t, "[Bash Defend Defend Strike Strike Strike Strike Defend Defend Strike]", draw.String())

	action = &StartBattleAction{
		Deck: nil,
		Draw: draw,
	}

	go Execute([]Action{action}, output, errorc)
	time.Sleep(time.Millisecond * 10)

	error := <-errorc
	assert.Equal(t, cards.ErrPileIsNilOrEmpty, error)
}
