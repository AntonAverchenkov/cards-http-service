package game

import (
	"fmt"
	"strings"
)

// Value represents one of the 13 card face values (Ace - King)
type Value uint8

// Value values
const (
	ValueAce Value = iota
	ValueTwo
	ValueThree
	ValueFour
	ValueFive
	ValueSix
	ValueSeven
	ValueEight
	ValueNine
	ValueTen
	ValueJack
	ValueQueen
	ValueKing
	ValuesTotalCount // a marker for the end of this enum
)

// ParseValue will parse the given card value string in short (a) or long (ace) forms
func ParseValue(str string) (Value, error) {
	switch strings.ToLower(str) {
	case "a", "ace":
		return ValueAce, nil
	case "2", "two":
		return ValueTwo, nil
	case "3", "three":
		return ValueThree, nil
	case "4", "four":
		return ValueFour, nil
	case "5", "five":
		return ValueFive, nil
	case "6", "six":
		return ValueSix, nil
	case "7", "seven":
		return ValueSeven, nil
	case "8", "eight":
		return ValueEight, nil
	case "9", "nine":
		return ValueNine, nil
	case "t", "ten":
		return ValueTen, nil
	case "j", "jack":
		return ValueJack, nil
	case "q", "queen":
		return ValueQueen, nil
	case "k", "king":
		return ValueKing, nil
	default:
		return 0, fmt.Errorf("could not parse '%s' as card value", str)
	}
}

func (v Value) String() string {
	return [...]string{
		"ace",
		"two",
		"three",
		"four",
		"five",
		"six",
		"seven",
		"eight",
		"nine",
		"ten",
		"jack",
		"queen",
		"king",
	}[v]
}

func (v Value) ShortString() string {
	return [...]string{
		"a",
		"2",
		"3",
		"4",
		"5",
		"6",
		"7",
		"8",
		"9",
		"t",
		"j",
		"q",
		"k",
	}[v]
}
