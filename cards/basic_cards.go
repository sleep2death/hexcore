package cards

// CardStrike is a Basic card. The Ironclad begins a run with 5 copies of Strike in the deck.
type CardStrike struct {
	*CardBase
}

// CreateCardStrike -
func CreateCardStrike() Card {
	actions := &Actions{
		Play: []Action{DealDamage},
	}

	return &CardStrike{
		CardBase: &CardBase{
			info:    &info{"strike", Attack, Red, Basic},
			base:    &Attrs{Damage: 6, Target: Enemy, Cost: 1, Actions: actions},
			upgrade: &Attrs{Damage: 3},
		},
	}
}

// CardDefend is a Basic card. The Ironclad begins a run with 4 copies of Defend in the deck.
type CardDefend struct {
	*CardBase
}

// CreateCardDefend -
func CreateCardDefend() Card {
	actions := &Actions{
		Play: []Action{GainBlock},
	}
	return &CardBash{
		CardBase: &CardBase{
			info:    &info{"defend", Skill, Red, Basic},
			base:    &Attrs{Block: 8, Target: Enemy, Cost: 2, Actions: actions},
			upgrade: &Attrs{Block: 3, Target: Enemy, Cost: 0},
		},
	}
}

// CardBash is an Attack card for the Ironclad. It deals 8 damage and applies Vulnerable for 2 turns. As a Basic card, the Ironclad always starts with one in his deck.
// Upon upgrade, Bash gains 2 damage and applies 1 more Vulnerable.
type CardBash struct {
	*CardBase
}

// CreateCardBash -
func CreateCardBash() Card {
	actions := &Actions{
		Play: []Action{Vulnerable},
	}
	return &CardBash{
		CardBase: &CardBase{
			info: &info{"bash", Attack, Red, Basic},
			base: &Attrs{Damage: 8, Target: Enemy, Cost: 2, Actions: actions},
		},
	}
}
