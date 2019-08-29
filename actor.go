package hexcore

// Actor of the battle
type Actor struct {
	id    string
	HP    uint
	MaxHP uint
}

// ID of the actor
func (a *Actor) ID() string {
	return a.id
}

// Player -
type Player struct {
	Actor
}

// Monster -
type Monster struct {
	Actor
}
