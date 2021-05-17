package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeckShuffle(t *testing.T) {
	deck := NewDeck()

	// back up the original deck
	original := make([]Card, len(deck.Cards))
	copy(original, deck.Cards)

	// shuffle the deck
	deck.Shuffle()

	// ensure the deck size has not changed
	require.Equal(t, len(original), deck.Len())

	// the cards must be the same but in different postions
	assert.ElementsMatch(t, original, deck.Cards)
	assert.NotEqual(t, original, deck.Cards)
}

func TestDeckDeal(t *testing.T) {
	deck := NewDeck()

	// back up the original deck
	original := make([]Card, len(deck.Cards))
	copy(original, deck.Cards)

	// deal one card
	card, err := deck.DealCard()
	require.NoError(t, err)
	assert.Equal(t, original[0], card)
	assert.Equal(t, len(original)-1, deck.Len())

	// deal the remaining cards
	var remaining []Card

	for deck.Len() != 0 {
		card, err := deck.DealCard()
		require.NoError(t, err)
		remaining = append(remaining, card)
	}

	assert.Equal(t, original[1:], remaining)

	// the deck is empty; try to deal one more
	_, err = deck.DealCard()
	assert.Error(t, err)
}

func TestDeckReturn(t *testing.T) {
	deck := NewDeck()

	// back up the original deck
	original := make([]Card, len(deck.Cards))
	copy(original, deck.Cards)

	// the deck is full, try to return one card, expect an error
	require.Error(t, deck.ReturnCard(Card{Value: ValueQueen, Suit: SuitDiamonds}))

	// deal 5 cards & return them
	var hand []Card

	for i := 0; i < 5; i++ {
		card, err := deck.DealCard()
		require.NoError(t, err)
		hand = append(hand, card)
	}

	for _, card := range hand {
		require.NoError(t, deck.ReturnCard(card))
	}

	// ensure the deck size has not changed
	require.Equal(t, len(original), deck.Len())

	// the 5 cards have been moved from the front of the deck to the back
	assert.Equal(t, original[5:], deck.Cards[:47])
	assert.Equal(t, original[:5], deck.Cards[47:])

	// return a card that already exists, expect an error
	_, err := deck.DealCard()
	require.NoError(t, err)
	require.Error(t, deck.ReturnCard(Card{Value: ValueAce, Suit: SuitSpades}))
}

func TestSerializeDeserialize(t *testing.T) {
	_, err := DeckDeserialize("ahq")
	require.Error(t, err)

	_, err = DeckDeserialize("ahqs3")
	require.Error(t, err)

	_, err = DeckDeserialize("ahqs3m")
	require.Error(t, err)

	expected := []Card{{
		Value: ValueAce,
		Suit:  SuitHearts,
	}, {
		Value: ValueQueen,
		Suit:  SuitSpades,
	}, {
		Value: ValueThree,
		Suit:  SuitDiamonds,
	}, {
		Value: ValueJack,
		Suit:  SuitSpades,
	}, {
		Value: ValueTen,
		Suit:  SuitClubs,
	}}

	deck, err := DeckDeserialize("ahqs3djstc")
	require.NoError(t, err)
	require.Equal(t, expected, deck.Cards)
	require.Equal(t, "ahqs3djstc", deck.Serialize())
}
