package poker

import "testing"

func TestParseCard(t *testing.T) {
	tests := []struct {
		input    string
		expected Card
	}{
		{"Ah", Card{Rank: Ace, Suit: Hearts}},
		{"10d", Card{Rank: Ten, Suit: Diamonds}},
		{"ks", Card{Rank: King, Suit: Spades}},
	}

	for _, tc := range tests {
		card, err := ParseCard(tc.input)
		if err != nil {
			t.Fatalf("unexpected error for %s: %v", tc.input, err)
		}
		if card != tc.expected {
			t.Fatalf("expected %v, got %v", tc.expected, card)
		}
	}
}

func TestParseCardInvalid(t *testing.T) {
	if _, err := ParseCard("1h"); err == nil {
		t.Fatal("expected error for invalid rank")
	}
	if _, err := ParseCard("Ahh"); err == nil {
		t.Fatal("expected error for invalid suit")
	}
}
