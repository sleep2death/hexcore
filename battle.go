package hexcore

import (

	// "github.com/sleep2death/hexcore/actors"

	"github.com/sleep2death/hexcore/cards"
)

var seed int64 = 9012

// Battle holds all card piles
type Battle struct {
	cards *cards.Manager
}
