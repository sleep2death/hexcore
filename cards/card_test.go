package cards

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCardCreate(t *testing.T) {
	card := CreateCardFunc["Strike"]()
	assert.EqualValues(t, card.Info(), "&{ID:strike Cost:1 Type:attack Color:red Rarity:basic Target:enemy}")
}
