package actors

// Actor of the battle
type Actor struct {
	HP    uint
	MaxHP uint
}

// Player -
type Player struct {
	Actor
}

// MonsterID -
type MonsterID uint

// Monster -
type Monster struct {
	Actor
}
