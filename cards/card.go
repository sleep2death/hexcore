package cards

import (
	"errors"
	"math/rand"
	"strings"
	"sync"

	"github.com/lithammer/shortuuid"
)

var (
	// ErrDrawNumber -
	ErrDrawNumber = errors.New("draw number should be larger than 0")

	// ErrDrawIndex -
	ErrDrawIndex = errors.New("draw index should be larger than 0 and less than len(cards) - 1")

	// ErrNotEnoughCards -
	ErrNotEnoughCards = errors.New("not enough card(s) to draw")

	// ErrCardNotExist -
	ErrCardNotExist = errors.New("card doesn't exist")

	// ErrPileIsNilOrEmpty -
	ErrPileIsNilOrEmpty = errors.New("target pile is nil or empty")
)

// Target is the target type of the card
type Target uint

const (
	// Enemy as the card target
	Enemy Target = iota + 1
	// AllEnemies as the card targets
	AllEnemies
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

// Status hold all changeable attributes of the card
type Status struct {
	Cost   int
	Damage int
	Block  int
	Heal   int
	Level  int
	Target Target
}

// Copy the status
func (n *Status) Copy() (s *Status) {
	s = &Status{
		Cost:   n.Cost,
		Damage: n.Damage,
		Block:  n.Block,
		Heal:   n.Heal,
		Target: n.Target,
		Level:  n.Level,
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

	return
}

// Card -
type Card struct {
	id    string
	name  string
	ctype CType

	base    *Status
	upgrade *Status
	current *Status
}

// Copy the card
func (card *Card) Copy() *Card {
	return &Card{
		// generate a new uuid for the card
		id: shortuuid.New(),
		// card name
		name: card.name,
		// card type
		ctype: card.ctype,
		// some cards may change the base status permantly in the battle
		// like card [Ritual Dagger] -  if this card kills an enemy then permanently increase this card's damage by 3(5)
		// if card["feed"] upgraded in the battle, then original card in the deck will also be upgraded
		// battle manager can use "base" status permanently change the card
		base: card.base,
		// card upgrade status will never be modified after created, so use a pointer is fine
		// when upgrade a card, use the base status add the upgrade status, then return the new upgraded status
		upgrade: card.upgrade,
		// card current status in battle
		current: card.base.Copy(),
	}
}

func (card *Card) String() string {
	return card.name
}

// Name of the card
func (card *Card) Name() string {
	return card.name
}

// ID return the uuid of the card
func (card *Card) ID() string {
	return card.id
}

// CType return the type of the card
func (card *Card) CType() CType {
	return card.ctype
}

// Base return the base status of the card
func (card *Card) Base() *Status {
	return card.base
}

// Upgrade return the upgrade of the card
func (card *Card) Upgrade() *Status {
	return card.upgrade
}

// UpgradeBase - upgrade card base status, also tit will upgrade the current status
func (card *Card) UpgradeBase() {
	card.base = card.base.Upgrade(card.upgrade)
	card.current = card.current.Upgrade(card.upgrade)
}

// UpgradeCurrent - upgrade card current status
func (card *Card) UpgradeCurrent() {
	card.current = card.current.Upgrade(card.upgrade)
}

// Current return the current status of the card
func (card *Card) Current() *Status {
	return card.current
}

// CreateCard -
func CreateCard(name string, ctype CType, base *Status, upgrade *Status) *Card {
	return &Card{
		// generate a new uuid for the card
		id: shortuuid.New(),
		// card name
		name: name,
		// card type
		ctype: ctype,
		// some cards may change the base status permantly in the battle
		// like card [Ritual Dagger] -  if this card kills an enemy then permanently increase this card's damage by 3(5)
		// if card["feed"] upgraded in the battle, then original card in the deck will also be upgraded
		// battle manager can use "base" status permanently change the card
		base: base,
		// card upgrade status will never be modified after created, so use a pointer is fine
		// when upgrade a card, use the base status add the upgrade status, then return the new upgraded status
		upgrade: upgrade,
		// card current status in battle
		current: base.Copy(),
	}
}

// Pile of cards
type Pile struct {
	cards []*Card
	// Lock
	mux sync.Mutex
}

// PileName -
type PileName int

const (
	// Deck Pile
	Deck PileName = iota
	// Draw Pile
	Draw
	// Hand Pile
	Hand
	// Discard Pile
	Discard
	// Exaust Pile
	Exaust
)

// String -
func (p *Pile) String() string {
	p.mux.Lock()
	defer p.mux.Unlock()

	s := ""
	for _, card := range p.cards {
		s += card.String() + " "
	}

	return "[" + strings.TrimSpace(s) + "]"
}

// Num - get the card number of the pile
func (p *Pile) Num() int {
	p.mux.Lock()
	defer p.mux.Unlock()
	// l := len(p.cards)
	return len(p.cards)
}

// Clear the pile
func (p *Pile) Clear() {
	p.mux.Lock()
	defer p.mux.Unlock()

	p.cards = nil
}

// AddToTop with the given card(s)
func (p *Pile) AddToTop(c ...*Card) {
	p.mux.Lock()
	defer p.mux.Unlock()
	p.cards = append(p.cards, c...)
}

// AddToBottom with the given card(s)
func (p *Pile) AddToBottom(c ...*Card) {
	p.mux.Lock()
	defer p.mux.Unlock()
	p.cards = append(c, p.cards...)
}

// Draw n card(s) from the source pile
func (p *Pile) Draw(source *Pile, n int) error {
	if n <= 0 {
		return ErrDrawNumber
	}
	if n > source.Num() {
		return ErrNotEnoughCards
	}

	p.mux.Lock()
	idx := len(source.cards) - n
	p.cards = append(p.cards, source.cards[idx:]...)
	p.mux.Unlock()

	source.mux.Lock()
	source.cards = source.cards[:idx]
	source.mux.Unlock()
	return nil
}

// Pick one card from source pile and add it to the top of the pile
func (p *Pile) Pick(source *Pile, id string) error {
	card, idx, err := source.FindCard(id)
	if err != nil {
		return err
	}

	p.AddToTop(card)
	source.RemoveCard(idx)
	return nil
}

// RemoveCard from the pile
func (p *Pile) RemoveCard(i int) (*Card, error) {
	p.mux.Lock()
	defer p.mux.Unlock()

	if i < 0 || i > len(p.cards)-1 {
		return nil, ErrDrawIndex
	}
	card := p.cards[i]
	copy(p.cards[i:], p.cards[i+1:])
	p.cards = p.cards[:len(p.cards)-1]
	return card, nil
}

// GetCard from the pile
func (p *Pile) GetCard(i int) (*Card, error) {
	p.mux.Lock()
	defer p.mux.Unlock()
	if i < 0 || i > len(p.cards)-1 {
		return nil, ErrDrawIndex
	}
	card := p.cards[i]
	return card, nil
}

// FindCard return both the card and card index of the pile
func (p *Pile) FindCard(id string) (card *Card, idx int, err error) {
	p.mux.Lock()
	defer p.mux.Unlock()
	if len(p.cards) == 0 {
		return nil, -1, ErrPileIsNilOrEmpty
	}

	for i, c := range p.cards {
		if c.ID() == id {
			return c, i, nil
		}
	}

	return nil, -1, ErrCardNotExist
}

// Shuffle the pile
func (p *Pile) Shuffle(seed *rand.Rand) {
	p.mux.Lock()
	defer p.mux.Unlock()

	if len(p.cards) == 0 {
		return
	}

	seed.Shuffle(len(p.cards), func(i, j int) { p.cards[i], p.cards[j] = p.cards[j], p.cards[i] })
}

// CopyCardsFrom -
func (p *Pile) CopyCardsFrom(source *Pile) error {
	if source == nil || source.Num() == 0 {
		return ErrPileIsNilOrEmpty
	}

	source.mux.Lock()
	for _, card := range source.cards {
		p.AddToTop(card.Copy())
	}
	source.mux.Unlock()

	return nil
}

// CreatePile by given seed and cardset
func CreatePile(cardSet []string) (p *Pile, err error) {
	p = &Pile{}

	if cardSet == nil || len(cardSet) == 0 {
		return p, nil
	}

	for _, s := range cardSet {
		if CreateCardFunc[s] == nil {
			// clear all the items reference by setting the slice to nil
			// see: https://stackoverflow.com/questions/16971741/how-do-you-clear-a-slice-in-go
			p.cards = nil
			return nil, ErrCardNotExist
		}

		card := CreateCardFunc[s]()
		p.AddToTop(card)
	}
	return
}

// CreateCardFunc map for generating cards
var CreateCardFunc = map[string](func() *Card){
	"Strike": CreateCardStrike,
	"Bash":   CreateCardBash,
	"Defend": CreateCardDefend,
}
