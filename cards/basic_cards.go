package cards

// CardStrike is a Basic card. The Ironclad begins a run with 5 copies of Strike in the deck.
type CardStrike struct {
	*CardBase
}

// CreateCardStrike -
func CreateCardStrike() (card Card) {
	card = &CardStrike{
		CardBase: &CardBase{
			info:    &info{"strike", Attack, Red, Basic},
			base:    &Status{Damage: 6, Target: Enemy, Cost: 1},
			upgrade: &Status{Damage: 3},
		},
	}
	card.Init()
	return
}

// CardDefend is a Basic card. The Ironclad begins a run with 4 copies of Defend in the deck.
type CardDefend struct {
	*CardBase
}

// CreateCardDefend -
func CreateCardDefend() (card Card) {
	card = &CardDefend{
		CardBase: &CardBase{
			info:    &info{"defend", Skill, Red, Basic},
			base:    &Status{Block: 5, Target: Enemy, Cost: 2},
			upgrade: &Status{Block: 3},
		},
	}
	card.Init()
	return
}

// CardBash is an Attack card for the Ironclad. It deals 8 damage and applies Vulnerable for 2 turns. As a Basic card, the Ironclad always starts with one in his deck.
// Upon upgrade, Bash gains 2 damage and applies 1 more Vulnerable.
type CardBash struct {
	*CardBase
}

// CreateCardBash -
func CreateCardBash() (card Card) {
	card = &CardBash{
		CardBase: &CardBase{
			info:    &info{"bash", Attack, Red, Basic},
			base:    &Status{Damage: 8, Target: Enemy, Cost: 2},
			upgrade: &Status{Damage: 2},
		},
	}
	card.Init()
	return
}
