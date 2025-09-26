package poker

import "testing"

func TestEvaluateBestHandCategories(t *testing.T) {
	tests := []struct {
		name     string
		cards    []string
		expected HandCategory
		ranks    []Rank
	}{
		{
			name:     "straight flush",
			cards:    []string{"Ah", "Kh", "Qh", "Jh", "Th", "2d", "3c"},
			expected: StraightFlush,
			ranks:    []Rank{Ace},
		},
		{
			name:     "four of a kind",
			cards:    []string{"Ah", "Ad", "Ac", "As", "Kd", "Qh", "2c"},
			expected: FourOfAKind,
			ranks:    []Rank{Ace, King},
		},
		{
			name:     "full house",
			cards:    []string{"Ah", "Ad", "As", "Kd", "Kc", "Qh", "2s"},
			expected: FullHouse,
			ranks:    []Rank{Ace, King},
		},
		{
			name:     "flush",
			cards:    []string{"Ah", "9h", "7h", "5h", "2h", "Kd", "Qc"},
			expected: Flush,
			ranks:    []Rank{Ace, Nine, Seven, Five, Two},
		},
		{
			name:     "straight",
			cards:    []string{"Ah", "Kd", "Qc", "Js", "Td", "2h", "3c"},
			expected: Straight,
			ranks:    []Rank{Ace},
		},
		{
			name:     "three of a kind",
			cards:    []string{"Ah", "Ad", "As", "Kd", "Qc", "5h", "2d"},
			expected: ThreeOfAKind,
			ranks:    []Rank{Ace, King, Queen},
		},
		{
			name:     "two pair",
			cards:    []string{"Ah", "Ad", "Kc", "Ks", "Qd", "3h", "2s"},
			expected: TwoPair,
			ranks:    []Rank{Ace, King, Queen},
		},
		{
			name:     "one pair",
			cards:    []string{"Ah", "Ad", "Kc", "Qd", "Jh", "3s", "2c"},
			expected: OnePair,
			ranks:    []Rank{Ace, King, Queen, Jack},
		},
		{
			name:     "high card",
			cards:    []string{"Ah", "Kd", "Qc", "Jd", "9s", "3h", "2c"},
			expected: HighCard,
			ranks:    []Rank{Ace, King, Queen, Jack, Nine},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cards := make([]Card, 0, len(tc.cards))
			for _, s := range tc.cards {
				cards = append(cards, MustParseCard(s))
			}

			rank, err := EvaluateBestHand(cards)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if rank.Category != tc.expected {
				t.Fatalf("expected category %v, got %v", tc.expected, rank.Category)
			}

			for i, expectedRank := range tc.ranks {
				if rank.Values[i] != expectedRank {
					t.Fatalf("expected value[%d]=%v, got %v", i, expectedRank, rank.Values[i])
				}
			}
		})
	}
}

func TestHandRankCompare(t *testing.T) {
	straight := handRankFromSlice(Straight, []Rank{Queen})
	flush := handRankFromSlice(Flush, []Rank{Jack, Nine, Seven, Five, Three})

	if straight.Compare(flush) != -1 {
		t.Fatalf("expected straight < flush")
	}

	a := handRankFromSlice(OnePair, []Rank{King, Queen, Jack, Ten})
	b := handRankFromSlice(OnePair, []Rank{King, Queen, Jack, Nine})

	if a.Compare(b) <= 0 {
		t.Fatalf("expected a > b based on kicker")
	}
}
