package cards

// CardStrike is a Basic card that deals 6 damage and is available to all characters.
type CardStrike struct {
	*CardBase
}

// CreateCardStrike -
func CreateCardStrike() Card {
	return &CardStrike{
		CardBase: &CardBase{
			info: &info{"strike", Attack, Red, Basic},
			base: &CardNum{Damage: 6, Target: Enemy, Cost: 1},
			actions: &actions{
				play: []CardAction{DealDamage},
			},
		},
	}
}
