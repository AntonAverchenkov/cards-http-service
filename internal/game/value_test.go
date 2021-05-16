package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseValue(t *testing.T) {
	successCases := []struct {
		str      string
		expected Value
	}{{
		str:      "five",
		expected: ValueFive,
	}, {
		str:      "six",
		expected: ValueSix,
	}, {
		str:      "9",
		expected: ValueNine,
	}, {
		str:      "a",
		expected: ValueAce,
	}}

	for _, test := range successCases {
		v, err := ParseValue(test.str)
		assert.NoError(t, err)
		assert.Equal(t, test.expected, v)
	}
}

func TestParseValueErrors(t *testing.T) {
	failureCases := []string{
		"0",
		"11",
		"111",
		"1111",
		"not a value",
		"ace of diamonds",
	}

	for _, test := range failureCases {
		_, err := ParseValue(test)
		assert.Error(t, err)
	}
}

func TestValueString(t *testing.T) {
	assert.Equal(t, "ace", ValueAce.String())
	assert.Equal(t, "three", ValueThree.String())
	assert.Equal(t, "jack", ValueJack.String())
	assert.Equal(t, "king", ValueKing.String())
}

func TestValueShortString(t *testing.T) {
	assert.Equal(t, "a", ValueAce.ShortString())
	assert.Equal(t, "3", ValueThree.ShortString())
	assert.Equal(t, "j", ValueJack.ShortString())
	assert.Equal(t, "k", ValueKing.ShortString())
}
