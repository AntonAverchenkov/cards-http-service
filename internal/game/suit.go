package game

import (
	"fmt"
	"strings"
)

// Suit represents one of the 4 card categories (Clubs / Hearts / Diamonds / Spades)
type Suit uint8

// Suit values
const (
	SuitClubs Suit = iota
	SuitHearts
	SuitDiamonds
	SuitSpades
	SuitsTotalCount // the total number of suits
)

// ParseSuit will parse a suit string in short (c) or long (clubs) forms
func ParseSuit(str string) (Suit, error) {
	switch strings.ToLower(str) {

	case "c", "clubs":
		return SuitClubs, nil
	case "h", "hearts":
		return SuitHearts, nil
	case "d", "diamonds":
		return SuitDiamonds, nil
	case "s", "spades":
		return SuitSpades, nil
	default:
		return 0, fmt.Errorf("could not parse '%s' as suit", str)
	}
}

// String representation of the suit
func (s Suit) String() string {
	return [...]string{
		"clubs",
		"hearts",
		"diamonds",
		"spades",
	}[s]
}

// Short string representation of the suit, used for persistence
func (s Suit) ShortString() string {
	return [...]string{
		"c",
		"h",
		"d",
		"s",
	}[s]
}
