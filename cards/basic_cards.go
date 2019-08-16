package cards

// CardStrike is a Basic card. The Ironclad begins a run with 5 copies of Strike in the deck.
type CardStrike struct {
	*CardBase
}

// CreateCardStrike -
func CreateCardStrike() Card {
	return &CardStrike{
		CardBase: &CardBase{
			info:    &info{"strike", Attack, Red, Basic},
			base:    &Numbers{Damage: 6, Target: Enemy, Cost: 1},
			upgrade: &Numbers{Damage: 3, Target: Enemy, Cost: 0},
			actions: &Actions{
				Play: []Action{DealDamage},
			},
		},
	}
}

// CardDefend is a Basic card. The Ironclad begins a run with 4 copies of Defend in the deck.
type CardDefend struct {
	*CardBase
}

// CreateCardDefend -
func CreateCardDefend() Card {
	return &CardBash{
		CardBase: &CardBase{
			info: &info{"defend", Skill, Red, Basic},
			base: &Numbers{Damage: 8, Target: Enemy, Cost: 2},
			actions: &Actions{
				Play: []Action{GainBlock},
			},
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
	return &CardBash{
		CardBase: &CardBase{
			info: &info{"bash", Attack, Red, Basic},
			base: &Numbers{Damage: 8, Target: Enemy, Cost: 2},
			actions: &Actions{
				Play: []Action{DealDamage, Vulnerable},
			},
		},
	}
}
