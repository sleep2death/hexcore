package cards

import (
	"errors"
	"fmt"
	"math/rand"

	"github.com/lithammer/shortuuid"
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
	// StatusT card type
	// Unplayable cards added to the deck during combat encounters. They are designed to bloat the deck and prevent the player from drawing beneficial cards, with some of them having additional negative effects. Unlike Curses, Status cards are removed from the deck at the end of combat.
	StatusT
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
	// SelectEnemy - get the enemy as the target of the card, waiting for user input
	SelectEnemy Action = iota + 1
	// SelectAllEnemies - get all enemies as the target of the card
	SelectAllEnemies
	// SelectSelf - get self as the target of the card
	SelectSelf
	// DealDamage to enemy(s)
	DealDamage
	// GainBlock for self
	GainBlock
	// Vulnerable creatures take 50% more damage from Attacks.
	Vulnerable
)

func (c Action) String() string {
	return [...]string{"nil", "select_enemy", "select_all_enemies", "select_self", "deal_damage", "gain_block", "vulnerable"}[c]
}

// info - basic infomation of the card
type info struct {
	ID     string
	CType  CType
	Color  Color
	Rarity Rarity
}

// Status hold all changeable attributes of the card
type Status struct {
	Cost    int
	Damage  int
	Block   int
	Heal    int
	Level   int
	Target  Target
	Actions *Actions
}

// Copy the status
func (n *Status) Copy() (s *Status) {
	s = &Status{
		Cost:    n.Cost,
		Damage:  n.Damage,
		Block:   n.Block,
		Heal:    n.Heal,
		Target:  n.Target,
		Level:   n.Level,
		Actions: n.Actions.Copy(),
	}
	return
}

// Upgrade a numbers with another numbers
func (n *Status) Upgrade(u *Status) (s *Status) {
	s = &Status{
		Cost:   n.Cost + u.Cost,
		Damage: n.Damage + u.Damage,
		Block:  n.Block + u.Block,
		Heal:   n.Heal + u.Heal,
		Target: n.Target + u.Target, // upgdate to the target number; if the u.Target == 0, then target not change
		Level:  n.Level + 1,         // update level + 1
	}

	if u.Actions != nil {
		s.Actions = u.Actions
	}

	return
}

// Actions - action holder of the card
type Actions struct {
	PreBattle  []Action
	PostBattle []Action
	PreTurn    []Action
	PostTurn   []Action
	Select     []Action
	Play       []Action
}

// Copy the actions
func (a *Actions) Copy() *Actions {
	c := &Actions{}

	if len(a.PreBattle) > 0 {
		c.PreBattle = make([]Action, len(a.PreBattle))
		copy(c.PreBattle, a.PreBattle)
	}

	if len(a.PostBattle) > 0 {
		c.PostBattle = make([]Action, len(a.PostBattle))
		copy(c.PostBattle, a.PostBattle)
	}

	if len(a.PreTurn) > 0 {
		c.PreTurn = make([]Action, len(a.PreTurn))
		copy(c.PreTurn, a.PreTurn)
	}

	if len(a.PostTurn) > 0 {
		c.PostTurn = make([]Action, len(a.PostTurn))
		copy(c.PostTurn, a.PostTurn)
	}

	if len(a.Select) > 0 {
		c.Select = make([]Action, len(a.Select))
		copy(c.Select, a.Select)
	}

	if len(a.Play) > 0 {
		c.Play = make([]Action, len(a.Play))
		copy(c.Play, a.Play)
	}

	return c
}

// CardBase -
type CardBase struct {
	*info
	id      string
	base    *Status
	upgrade *Status
	current *Status
}

// Copy the card
func (card *CardBase) Copy() Card {
	c := &CardBase{
		// generate a new uuid for the card
		id: shortuuid.New(),
		// card info will never be modified after created, so use a pointer is fine
		info: card.info,
		// card upgrade status will never be modified after created, so use a pointer is fine
		// when upgrade a card, use the base status add the upgrade status, then return the new upgraded status
		upgrade: card.upgrade,
		// some cards may change the base status permantly in the battle
		// like card [Ritual Dagger] -  if this card kills an enemy then permanently increase this card's damage by 3(5)
		// if card["feed"] upgraded in the battle, then original card in the deck will also be upgraded
		// manager can use "base" status permanently change the card
		base: card.base,
		// card current status in battle
	}

	if card.current != nil {
		c.current = card.current.Copy()
	}

	return c
}

// Init the card by copying the base status to current
func (card *CardBase) Init() error {
	if len(card.id) > 0 || card.current != nil {
		return fmt.Errorf("card %v has been initialized already", card)
	}
	card.id = shortuuid.New()       //generate a new uuid for the card
	card.current = card.base.Copy() // copy the base status to the current status

	return nil
}

// Info return the basic information of the card
func (card *CardBase) String() string {
	return fmt.Sprintf("[%s]", card.info.ID)
}

// ID return the uuid of the card
func (card *CardBase) ID() string {
	return card.id
}

// Base return the base numbers of the card
func (card *CardBase) Base() *Status {
	return card.base
}

// Current return the current numbers of the card
func (card *CardBase) Current() *Status {
	return card.current
}

// Upgrade the card
func (card *CardBase) Upgrade() {
	card.base = card.base.Upgrade(card.upgrade)
}

// Card interface
type Card interface {
	String() string
	ID() string

	Base() *Status
	Current() *Status

	Copy() Card
	Upgrade()
	Init() error
}

// Pile of cards
type Pile struct {
	seed  int64
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

// FindCardByID return the card index with given id
func (p *Pile) FindCardByID(id string) int {
	for i, c := range p.cards {
		if c.ID() == id {
			return i
		}
	}

	return -1
}

// Shuffle the pile
func (p *Pile) Shuffle() {
	rand.Seed(p.seed)
	rand.Shuffle(len(p.cards), func(i, j int) { p.cards[i], p.cards[j] = p.cards[j], p.cards[i] })
}

// CardsNum - get the card number of the pile
func (p *Pile) CardsNum() int {
	return len(p.cards)
}

// CreateCardByName - create the card by the given name
func (p *Pile) CreateCardByName(cardSet []string) error {
	for _, s := range cardSet {
		if CreateCardFunc[s] == nil {
			// clear all the items by setting the slice to nil
			// see: https://stackoverflow.com/questions/16971741/how-do-you-clear-a-slice-in-go
			p.cards = nil
			return fmt.Errorf("create function for card [%s] not found", s)
		}

		card := CreateCardFunc[s]()
		p.AddToTop(card)
	}
	return nil
}

// CreateCardFunc map for generating cards
var CreateCardFunc = map[string](func() Card){
	"Strike": CreateCardStrike,
	"Bash":   CreateCardBash,
	"Defend": CreateCardDefend,
}
