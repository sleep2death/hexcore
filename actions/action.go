package actions

import (
	"github.com/sleep2death/hexcore/actors"
)

// Action -
type Action interface {
	// Exec - execute the action
	Exec(p actors.Player, m []actors.Monster) error
}

// DealDamage action
type DealDamage struct {
	amount uint
}

// Exec -
func (a *DealDamage) Exec(p actors.Player, m []actors.Monster) error {
	return nil
}
