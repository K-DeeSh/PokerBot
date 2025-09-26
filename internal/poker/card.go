package poker

import (
	"fmt"
	"strings"
)

// Suit enumerates the four card suits.
type Suit int

const (
	Clubs Suit = iota
	Diamonds
	Hearts
	Spades
)

var suitToString = map[Suit]string{
	Clubs:    "c",
	Diamonds: "d",
	Hearts:   "h",
	Spades:   "s",
}

var stringToSuit = map[string]Suit{
	"c": Clubs,
	"d": Diamonds,
	"h": Hearts,
	"s": Spades,
}

// Rank represents the numerical value of a card from Two to Ace.
type Rank int

const (
	Two Rank = iota
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
	Ace
)

var rankToString = map[Rank]string{
	Two:   "2",
	Three: "3",
	Four:  "4",
	Five:  "5",
	Six:   "6",
	Seven: "7",
	Eight: "8",
	Nine:  "9",
	Ten:   "T",
	Jack:  "J",
	Queen: "Q",
	King:  "K",
	Ace:   "A",
}

var stringToRank = map[string]Rank{
	"2":  Two,
	"3":  Three,
	"4":  Four,
	"5":  Five,
	"6":  Six,
	"7":  Seven,
	"8":  Eight,
	"9":  Nine,
	"T":  Ten,
	"10": Ten,
	"J":  Jack,
	"Q":  Queen,
	"K":  King,
	"A":  Ace,
}

// Card represents a single playing card.
type Card struct {
	Rank Rank
	Suit Suit
}

func (c Card) String() string {
	return rankToString[c.Rank] + suitToString[c.Suit]
}

// ParseCard parses a textual representation like "Ah" or "10d" into a Card.
func ParseCard(input string) (Card, error) {
	trimmed := strings.TrimSpace(strings.ToUpper(input))
	if len(trimmed) < 2 || len(trimmed) > 3 {
		return Card{}, fmt.Errorf("invalid card format: %s", input)
	}

	var rankStr, suitStr string
	if len(trimmed) == 3 {
		rankStr = trimmed[:2]
		suitStr = trimmed[2:]
	} else {
		rankStr = trimmed[:1]
		suitStr = trimmed[1:]
	}

	rank, ok := stringToRank[rankStr]
	if !ok {
		return Card{}, fmt.Errorf("invalid card rank: %s", rankStr)
	}

	suit, ok := stringToSuit[strings.ToLower(suitStr)]
	if !ok {
		return Card{}, fmt.Errorf("invalid card suit: %s", suitStr)
	}

	return Card{Rank: rank, Suit: suit}, nil
}

// MustParseCard is a helper that panics on invalid input; intended for tests.
func MustParseCard(input string) Card {
	card, err := ParseCard(input)
	if err != nil {
		panic(err)
	}
	return card
}

// AllCards returns the 52 cards in a standard deck.
func AllCards() []Card {
	cards := make([]Card, 0, 52)
	for suit := Clubs; suit <= Spades; suit++ {
		for rank := Two; rank <= Ace; rank++ {
			cards = append(cards, Card{Rank: rank, Suit: suit})
		}
	}
	return cards
}

// ContainsCard checks if the given slice already features the card.
func ContainsCard(cards []Card, target Card) bool {
	for _, c := range cards {
		if c == target {
			return true
		}
	}
	return false
}
