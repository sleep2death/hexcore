package cards

// CreateCardStrike -
func CreateCardStrike() (card *Card) {
	card = CreateCard("Strike", Attack, &Status{Damage: 6, Target: Enemy, Cost: 1}, &Status{Damage: 3})
	return
}

// CreateCardDefend -
func CreateCardDefend() (card *Card) {
	card = CreateCard("Defend", Skill, &Status{Block: 5, Target: Self, Cost: 1}, &Status{Block: 3})
	return
}

// CreateCardBash -
func CreateCardBash() (card *Card) {
	card = CreateCard("Bash", Attack, &Status{Damage: 8, Target: Enemy, Cost: 2}, &Status{Damage: 2})
	return
}
