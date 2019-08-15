package hexcore

import "fmt"

// CardColor defines the color of the card
type CardColor uint

const (
	// Red is for warrior cards
	Red CardColor = iota
	// Green is for roger cards
	Green
	// Blue is for wizard cards
	Blue
	// ColorLess is for  neutral cards (grey)
	ColorLess
	// CurseC is for curse cards (also grey)
	CurseC
)

// CardTarget is the target(s) of the card
type CardTarget uint

const (
	// Enemy as the card target
	Enemy CardTarget = iota
	// AllEnemy as the card targets
	AllEnemy
	// Self as the card targets
	Self
	// None target
	None
	// SelfAndEnemy as the card targets
	SelfAndEnemy
	// All as the card targets
	All
)

// CardRarity is the rarity of the card
type CardRarity uint

const (
	// Basic rarity
	Basic CardRarity = iota
	// Special rarity
	Special
	// Common rarity
	Common
	// UnCommon rarity
	UnCommon
	// Rare rarity
	Rare
	// CurseR rarity
	CurseR
)

// CardType is the type of the card
type CardType uint

const (
	// Attack card type
	Attack CardType = iota
	// Skill card type
	Skill
	// Power card type
	Power
	// Status card type
	Status
	// CurseT card type
	CurseT
)

// Card struct
type Card struct {
	ID             string
	Name           string
	ImgURL         string
	Cost           int
	RawDescription string
	Type           CardType
	Color          CardColor
	Rarity         CardRarity
	Target         CardTarget
}

// ToString return the general information of the card
func (card *Card) ToString() string {
	return fmt.Sprintf("Card: %s", card.ID)
}

// ICard interface
type ICard interface {
	ToString() string
}
