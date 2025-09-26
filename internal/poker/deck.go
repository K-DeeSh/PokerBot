package poker

import "math/rand"

// BuildDeck returns a deck excluding the specified cards.
func BuildDeck(excluded []Card) []Card {
	excludeMap := make(map[Card]struct{}, len(excluded))
	for _, c := range excluded {
		excludeMap[c] = struct{}{}
	}

	deck := make([]Card, 0, 52-len(excluded))
	for _, c := range AllCards() {
		if _, exists := excludeMap[c]; exists {
			continue
		}
		deck = append(deck, c)
	}
	return deck
}

// DrawCards removes cards from the deck slice and returns the drawn portion.
func DrawCards(deck *[]Card, n int, rng *rand.Rand) []Card {
	cards := *deck
	if n > len(cards) {
		n = len(cards)
	}

	// Fisher-Yates partial shuffle
	for i := 0; i < n; i++ {
		j := i + rng.Intn(len(cards)-i)
		cards[i], cards[j] = cards[j], cards[i]
	}

	drawn := append([]Card(nil), cards[:n]...)
	*deck = cards[n:]
	return drawn
}
