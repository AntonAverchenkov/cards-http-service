package game

import (
	"fmt"
	"strings"
)

type Card struct {
	Suit
	Value
}

// ParseCard will parse the given string into a card object from short ("ad") or long ("ace of diamonds") forms
func ParseCard(str string) (Card, error) {
	str = strings.ToLower(str)

	// long format
	tokens := strings.Split(str, " of ")

	if len(tokens) != 2 {
		// check if this is short format instead
		if len(str) != 2 {
			return Card{}, fmt.Errorf("unexpected card format '%s'", str)
		}

		// short fomrat
		tokens = []string{
			string(str[0]),
			string(str[1]),
		}
	}

	v, err := ParseValue(tokens[0])
	if err != nil {
		return Card{}, fmt.Errorf("couldn't parse card's first token: %w", err)
	}

	s, err := ParseSuit(tokens[1])
	if err != nil {
		return Card{}, fmt.Errorf("couldn't parse card's second token: %w", err)
	}

	return Card{
		Value: v,
		Suit:  s,
	}, nil
}

// String will return something like "ace of hearts"
func (c Card) String() string {
	return fmt.Sprintf("%s of %s", c.Value, c.Suit)
}

// ShortString will return something like "ah"
func (c Card) ShortString() string {
	return fmt.Sprintf("%s%s", c.Value.ShortString(), c.Suit.ShortString())
}
