package cards

import (
	"fmt"
)

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

func (c CardColor) String() string {
	return [...]string{"red", "green", "blue", "colorless", "curse"}[c]
}

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

func (c CardTarget) String() string {
	return [...]string{"enemy", "allEnemy", "self", "none", "selfAndEnemy", "all"}[c]
}

// CardRarity is the rarity of the card
type CardRarity uint

const (
	// Basic rarity
	// Basic cards are the default cards from the starting deck for your class. They have the same grey banner as Commons, though certain events treat them as a lower tier when offered or transformed.
	Basic CardRarity = iota
	// Special rarity
	// Special cards cannot be obtained through normal card-drops.
	Special
	// Common rarity
	// Common cards have a grey banner
	Common
	// Uncommon rarity
	// Uncommon cards have a blue banner
	Uncommon
	// Rare rarity
	// Rare cards have a yellow/gold banne
	Rare
	// CurseR rarity
	CurseR
)

func (c CardRarity) String() string {
	return [...]string{"basic", "special", "common", "uncommon", "rare", "curse"}[c]
}

// CardType is the type of the card
type CardType uint

const (
	// Attack card type
	// A reusable card (Unless it has Exhaust) that deals direct damage to an enemy and may have a secondary effect.
	Attack CardType = iota
	// Skill card type
	// A reusable card (Unless it has Exhaust) that has more unique effects to it. There isn't a clear direction with offensiveness and defensiveness unlike attacks.
	Skill
	// Power card type
	// A permanent upgrade for the entire combat encounter. Some Powers give flat stats like Strength or Dexterity. Others require certain conditions to be met that combat. Each copy of a given power can only be played once per combat.
	Power
	// Status card type
	// Unplayable cards added to the deck during combat encounters. They are designed to bloat the deck and prevent the player from drawing beneficial cards, with some of them having additional negative effects. Unlike Curses, Status cards are removed from the deck at the end of combat.
	Status
	// CurseT card type
	// Unplayable cards added to the deck during in-game events. Similar to status cards they are designed to bloat the deck and prevent the player from drawing beneficial cards, with some of them having additional negative effects during combat. Unlike Statuses, Curse cards persist in the players' deck until removed by other means.
	CurseT
)

func (c CardType) String() string {
	return [...]string{"attack", "skill", "power", "status", "curse"}[c]
}

// CardAction is the type of the card actions
type CardAction uint

const (
	// DealDamage to enemy(s)
	DealDamage CardAction = iota
	// GainBlock for self
	GainBlock
)

func (c CardAction) String() string {
	return [...]string{"deal_damage", "gain_block"}[c]
}

// info - basic infomation of the card
type info struct {
	ID     string
	CType  CardType
	Color  CardColor
	Rarity CardRarity
}

// CardNum hold all numbers of the card
type CardNum struct {
	Cost   int
	Damage int
	Block  int
	Heal   int
	Target CardTarget
}

// actions - action holder of the card
type actions struct {
	preBattle  []CardAction
	postBattle []CardAction
	preTurn    []CardAction
	postTurn   []CardAction
	play       []CardAction
}

func (a *actions) PreBattle() []CardAction {
	return a.preBattle
}

func (a *actions) PostBattle() []CardAction {
	return a.postBattle
}

func (a *actions) PreTurn() []CardAction {
	return a.preTurn
}

func (a *actions) PostTurn() []CardAction {
	return a.postTurn
}

func (a *actions) Play() []CardAction {
	return a.play
}

// CardBase -
type CardBase struct {
	*info
	*actions
	base    *CardNum
	current *CardNum
}

// Info return the basic information of the card
func (card *CardBase) Info() string {
	return fmt.Sprintf("%+v", card.info)
}

// GetBase return the base numbers of the card
func (card *CardBase) GetBase() *CardNum {
	return card.base
}

// GetCurrent return the current numbers of the card
func (card *CardBase) GetCurrent() *CardNum {
	return card.current
}

// Card interface
type Card interface {
	Info() string
	GetBase() *CardNum
	GetCurrent() *CardNum
	PreBattle() []CardAction
	PostBattle() []CardAction
	PreTurn() []CardAction
	PostTurn() []CardAction
	Play() []CardAction
}

// CreateCardFunc map for generating cards
var CreateCardFunc = map[string](func() Card){
	"Strike": CreateCardStrike,
}
