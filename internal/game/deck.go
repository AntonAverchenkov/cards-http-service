package game

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	DeckCapacity = int(SuitsTotalCount) * int(ValuesTotalCount)
)

type Deck struct {
	Cards []Card

	// rng is a random number generator used to shuffle this deck
	rng *rand.Rand
}

// NewDeck initializes the deck with 52 unique cards in sorted order
func NewDeck() *Deck {
	cards := make([]Card, 0, 52)

	for suit := SuitClubs; suit < SuitsTotalCount; suit++ {
		for value := ValueAce; value < ValuesTotalCount; value++ {
			cards = append(cards, Card{
				Value: value,
				Suit:  suit,
			})
		}
	}

	return &Deck{
		Cards: cards,
		rng:   rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// DealCard removes the top card from the deck (subtracts it from the slice's front) and returns it
func (d *Deck) DealCard() (Card, error) {
	if len(d.Cards) == 0 {
		return Card{}, fmt.Errorf("the deck is empty")
	}

	top := d.Cards[0]

	// remove the top card from the deck
	d.Cards = d.Cards[1:]

	return top, nil
}

// ReturnCard adds the given card to the deck (at the end of the slice)
func (d *Deck) ReturnCard(card Card) error {
	if len(d.Cards) >= DeckCapacity {
		return fmt.Errorf("the deck is full")
	}

	if d.find(card) != -1 {
		return fmt.Errorf("the card '%s' already exists in the deck", card)
	}

	d.Cards = append(d.Cards, card)

	return nil
}

// Shuffle permutes the deck of cards using a pseudo-random algorithm seeded at deck creation time
func (d *Deck) Shuffle() {
	// apparently go already has a standard library implementation :)
	d.rng.Shuffle(len(d.Cards), func(i, j int) {
		d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]
	})
}

// Len returns the current deck's size
func (d *Deck) Len() int {
	return len(d.Cards)
}

// find returns the index of the given card or -1 if it does not exist
func (d *Deck) find(card Card) int {
	for i, c := range d.Cards {
		if card == c {
			return i
		}
	}

	return -1
}
