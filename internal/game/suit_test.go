package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSuit(t *testing.T) {
	successCases := []struct {
		str      string
		expected Suit
	}{{
		str:      "spades",
		expected: SuitSpades,
	}, {
		str:      "Clubs",
		expected: SuitClubs,
	}, {
		str:      "d",
		expected: SuitDiamonds,
	}, {
		str:      "H",
		expected: SuitHearts,
	}}

	for _, test := range successCases {
		s, err := ParseSuit(test.str)
		assert.NoError(t, err)
		assert.Equal(t, test.expected, s)
	}
}

func TestParseSuitErrors(t *testing.T) {
	failureCases := []string{
		"",
		"a",
		"aa",
		"aaa",
		"aaaa",
		"not a suit",
	}

	for _, test := range failureCases {
		_, err := ParseSuit(test)
		assert.Error(t, err)
	}
}

func TestSuitString(t *testing.T) {
	assert.Equal(t, "spades", SuitSpades.String())
	assert.Equal(t, "clubs", SuitClubs.String())
	assert.Equal(t, "diamonds", SuitDiamonds.String())
	assert.Equal(t, "hearts", SuitHearts.String())
}

func TestSuitShortString(t *testing.T) {
	assert.Equal(t, "s", SuitSpades.ShortString())
	assert.Equal(t, "c", SuitClubs.ShortString())
	assert.Equal(t, "d", SuitDiamonds.ShortString())
	assert.Equal(t, "h", SuitHearts.ShortString())
}
