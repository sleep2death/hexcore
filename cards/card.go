package cards

import (
	"errors"
	"fmt"
	"math/rand"
)

// Color defines the color of the card
type Color uint

const (
	// Red is for warrior cards
	Red Color = iota
	// Green is for roger cards
	Green
	// Blue is for wizard cards
	Blue
	// ColorLess is for  neutral cards (grey)
	ColorLess
	// CurseC is for curse cards (also grey)
	CurseC
)

func (c Color) String() string {
	return [...]string{"red", "green", "blue", "colorless", "curse"}[c]
}

// Target is the target type of the card
type Target uint

const (
	// Enemy as the card target
	Enemy Target = iota + 1
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

func (c Target) String() string {
	return [...]string{"nil", "enemy", "allEnemy", "self", "none", "selfAndEnemy", "all"}[c]
}

// Rarity is the rarity of the card
type Rarity uint

const (
	// Basic rarity
	// Basic cards are the default cards from the starting deck for your class. They have the same grey banner as Commons, though certain events treat them as a lower tier when offered or transformed.
	Basic Rarity = iota
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

func (c Rarity) String() string {
	return [...]string{"basic", "special", "common", "uncommon", "rare", "curse"}[c]
}

// CType is the type of the card
type CType uint

const (
	// Attack card type
	// A reusable card (Unless it has Exhaust) that deals direct damage to an enemy and may have a secondary effect.
	Attack CType = iota
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

func (c CType) String() string {
	return [...]string{"attack", "skill", "power", "status", "curse"}[c]
}

// Action is the type of the card actions
type Action uint

const (
	// DealDamage to enemy(s)
	DealDamage Action = iota + 1
	// GainBlock for self
	GainBlock
	// Vulnerable creatures take 50% more damage from Attacks.
	Vulnerable
)

func (c Action) String() string {
	return [...]string{"nil", "deal_damage", "gain_block", "vulnerable"}[c]
}

// info - basic infomation of the card
type info struct {
	ID     string
	CType  CType
	Color  Color
	Rarity Rarity
}

// Attrs hold all changeable attributes of the card
type Attrs struct {
	Cost    int
	Damage  int
	Block   int
	Heal    int
	Level   int
	Target  Target
	Actions *Actions
}

// Upgrade a numbers with another numbers
func (n *Attrs) Upgrade(u *Attrs) (attr *Attrs) {
	attr = &Attrs{
		Cost:   n.Cost + u.Cost,
		Damage: n.Damage + u.Damage,
		Block:  n.Block + u.Block,
		Heal:   n.Heal + u.Heal,
		Target: n.Target + u.Target, // upgdate to the target number
		Level:  n.Level + 1,         // update level + 1
	}

	if u.Actions != nil {
		attr.Actions = u.Actions
	}

	return
}

// Actions - action holder of the card
type Actions struct {
	PreBattle  []Action
	PostBattle []Action
	PreTurn    []Action
	PostTurn   []Action
	Play       []Action
}

// CardBase -
type CardBase struct {
	*info
	base    *Attrs
	upgrade *Attrs
	current *Attrs
}

// Info return the basic information of the card
func (card *CardBase) String() string {
	return fmt.Sprintf("[%s]", card.info.ID)
}

// Base return the base numbers of the card
func (card *CardBase) Base() *Attrs {
	return card.base
}

// Current return the current numbers of the card
func (card *CardBase) Current() *Attrs {
	return card.current
}

// Upgrade the card
func (card *CardBase) Upgrade() {
	card.base = card.base.Upgrade(card.upgrade)
}

// Card interface
type Card interface {
	String() string
	Base() *Attrs
	Current() *Attrs
	Upgrade()
}

// Pile of cards
type Pile struct {
	cards []Card
}

// AddToTop with the given card(s)
func (p *Pile) AddToTop(c ...Card) {
	p.cards = append(p.cards, c...)
}

// AddToBottom with the given card(s)
func (p *Pile) AddToBottom(c ...Card) {
	p.cards = append(c, p.cards...)
}

// Draw n card(s) to the target pile
func (p *Pile) Draw(n int, target *Pile) error {
	if n <= 0 {
		return fmt.Errorf("n(%d) should be larger than 0", n)
	}
	if n > len(p.cards) {
		return errors.New("not enough card(s) to draw")
	}

	idx := len(p.cards) - n

	target.cards = append(target.cards, p.cards[idx:]...)
	p.cards = p.cards[:idx]
	return nil
}

// Shuffle the pile
func (p *Pile) Shuffle(seed int64) {
	rand.Seed(seed)
	rand.Shuffle(len(p.cards), func(i, j int) { p.cards[i], p.cards[j] = p.cards[j], p.cards[i] })
}

// CreateCardFunc map for generating cards
var CreateCardFunc = map[string](func() Card){
	"Strike": CreateCardStrike,
	"Bash":   CreateCardBash,
	"Defend": CreateCardDefend,
}
