package actors

import (
	"github.com/sleep2death/hexcore/cards"
)

// Creature -
type Creature struct {
	hp    uint
	hpMax uint
	block uint
}

// Player -
type Player struct {
	cards.Manager
	Creature
}
