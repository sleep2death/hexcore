package cards

var seed int64 = 9012

// Manager of all the cards in player's pocket
type Manager struct {
	deck    *Pile
	draw    *Pile
	hand    *Pile
	discard *Pile
	exaust  *Pile
}

// Init the card manager by filling the deck with given cards
func (m *Manager) Init(cardSet []string) error {
	m.deck = &Pile{seed: seed}
	return m.deck.CreateCardByName(cardSet)
}

// Shuffle all cards in the deck pile to draw pile
func (m *Manager) Shuffle() {
}

// ReShuffle all cards in the discard pile to draw pile when draw pile is empty
func (m *Manager) ReShuffle() {
}

// Draw cards from draw pile to hand pile
func (m *Manager) Draw() {
}

// Exaust cards from hand pile to exaust pile
func (m *Manager) Exaust(card Card) {
}

// Discard cards from hand pile to discard pile
func (m *Manager) Discard(card Card) {
}
