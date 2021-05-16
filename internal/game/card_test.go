package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCard(t *testing.T) {
	successCases := []struct {
		str      string
		expected Card
	}{{
		str:      "five of spades",
		expected: Card{Value: ValueFive, Suit: SuitSpades},
	}, {
		str:      "Six of Clubs",
		expected: Card{Value: ValueSix, Suit: SuitClubs},
	}, {
		str:      "9d",
		expected: Card{Value: ValueNine, Suit: SuitDiamonds},
	}, {
		str:      "ah",
		expected: Card{Value: ValueAce, Suit: SuitHearts},
	}}

	for _, test := range successCases {
		c, err := ParseCard(test.str)
		assert.NoError(t, err)
		assert.Equal(t, test.expected, c)
	}
}

func TestParseCardErrors(t *testing.T) {
	failureCases := []string{
		"",
		" ",
		"qw",
		"abcd",
		"abcde",
		"ace of",
		"queen of",
		"not a card",
		"king of something",
	}

	for _, test := range failureCases {
		_, err := ParseCard(test)
		assert.Error(t, err)
	}
}

func TestCardString(t *testing.T) {
	assert.Equal(
		t,
		"ace of spades",
		Card{Value: ValueAce, Suit: SuitSpades}.String(),
	)
	assert.Equal(
		t,
		"three of clubs",
		Card{Value: ValueThree, Suit: SuitClubs}.String(),
	)
	assert.Equal(
		t,
		"jack of diamonds",
		Card{Value: ValueJack, Suit: SuitDiamonds}.String(),
	)
	assert.Equal(
		t,
		"king of hearts",
		Card{Value: ValueKing, Suit: SuitHearts}.String(),
	)
}

func TestCardShortString(t *testing.T) {
	assert.Equal(t, "as", Card{Value: ValueAce, Suit: SuitSpades}.ShortString())
	assert.Equal(t, "3c", Card{Value: ValueThree, Suit: SuitClubs}.ShortString())
	assert.Equal(t, "jd", Card{Value: ValueJack, Suit: SuitDiamonds}.ShortString())
	assert.Equal(t, "kh", Card{Value: ValueKing, Suit: SuitHearts}.ShortString())
}
