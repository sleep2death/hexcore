package cards

// CardStrike is a Basic card that deals 6 damage and is available to all characters.
type CardStrike struct {
	*info
	*nums
}

// CreateCardStrike -
func CreateCardStrike() ICard {
	return &CardStrike{
		info: &info{"strike", 1, Attack, Red, Basic, Enemy},
		nums: &nums{damage: 6},
	}
}
