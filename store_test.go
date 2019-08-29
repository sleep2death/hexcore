package hexcore

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStore(t *testing.T) {
	a := &State{}
	a.SetNum(1)
	assert.Equal(t, 1, a.Num())

	b := &State{}
	b.SetNum(2)
	assert.Equal(t, 2, b.Num())

	c := &State{}
	c.SetNum(3)
	assert.Equal(t, 3, c.Num())

	// test if the id worksfine in goroutines
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		GetStore().AddState(a)
		wg.Done()
	}()

	go func() {
		GetStore().AddState(b)
		wg.Done()
	}()

	wg.Wait()
	id := GetStore().AddState(c)
	assert.Equal(t, 2, id)

	s := GetStore().State(0)
	assert.Equal(t, s, a)
}

func TestPiles(t *testing.T) {
	p := make(Pile, 0)
	for i := 0; i < 10; i++ {
		card := &TestCard{}
		card.SetID(strconv.Itoa(i))
		p = append(p, card)
	}

	s := &State{}
	s.SetPile(Deck, &p)

	h := make(Pile, 0)
	s.SetPile(Draw, &h)

	assert.Equal(t, "0", (*s.GetPile(Deck))[0].ID())
	assert.Equal(t, "9", (*s.GetPile(Deck))[9].ID())

	seed := rand.New(rand.NewSource(99))
	deck := s.GetPile(Deck)
	draw := s.GetPile(Draw)

	s.Shuffle(Deck, seed)
	assert.Equal(t, "&[<card 1> <card 7> <card 4> <card 0> <card 9> <card 2> <card 3> <card 5> <card 8> <card 6>]", fmt.Sprint(deck))

	s.Draw(Deck, Draw)
	assert.Equal(t, "&[<card 1> <card 7> <card 4> <card 0> <card 9> <card 2> <card 3> <card 5> <card 8>]", fmt.Sprint(deck))
	assert.Equal(t, "&[<card 6>]", fmt.Sprint(draw))

	s.Pick("7", Deck, Draw)
	assert.Equal(t, "&[<card 1> <card 4> <card 0> <card 9> <card 2> <card 3> <card 5> <card 8>]", fmt.Sprint(deck))
	assert.Equal(t, "&[<card 6> <card 7>]", fmt.Sprint(draw))

	s.Copy(Draw, Hand)
	assert.Equal(t, "&[<card copy:1 of <6>> <card copy:1 of <7>>]", fmt.Sprint(s.GetPile(Hand)))
}
