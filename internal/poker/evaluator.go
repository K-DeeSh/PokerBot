package poker

import (
	"errors"
	"sort"
)

// HandCategory enumerates poker hand categories ordered by strength.
type HandCategory int

const (
	HighCard HandCategory = iota
	OnePair
	TwoPair
	ThreeOfAKind
	Straight
	Flush
	FullHouse
	FourOfAKind
	StraightFlush
)

// HandRank captures the strength of a five-card poker hand.
type HandRank struct {
	Category HandCategory
	Values   [5]Rank
}

// Compare returns 1 if h > other, -1 if h < other, 0 otherwise.
func (h HandRank) Compare(other HandRank) int {
	if h.Category > other.Category {
		return 1
	}
	if h.Category < other.Category {
		return -1
	}

	for i := range h.Values {
		if h.Values[i] > other.Values[i] {
			return 1
		}
		if h.Values[i] < other.Values[i] {
			return -1
		}
	}
	return 0
}

// EvaluateBestHand finds the best five-card hand from the provided cards.
func EvaluateBestHand(cards []Card) (HandRank, error) {
	n := len(cards)
	if n < 5 {
		return HandRank{}, errors.New("at least five cards required")
	}

	best := HandRank{}
	hasBest := false

	for i := 0; i < n-4; i++ {
		for j := i + 1; j < n-3; j++ {
			for k := j + 1; k < n-2; k++ {
				for l := k + 1; l < n-1; l++ {
					for m := l + 1; m < n; m++ {
						hand := []Card{cards[i], cards[j], cards[k], cards[l], cards[m]}
						eval := evaluateFive(hand)
						if !hasBest || eval.Compare(best) > 0 {
							best = eval
							hasBest = true
						}
					}
				}
			}
		}
	}

	return best, nil
}

func evaluateFive(cards []Card) HandRank {
	suitCounts := make(map[Suit]int, 4)
	rankCounts := make(map[Rank]int, len(cards))
	rankMask := 0

	for _, c := range cards {
		suitCounts[c.Suit]++
		rankCounts[c.Rank]++
		rankMask |= 1 << int(c.Rank)
	}

	allRanksDesc := sortedRanksByCount(rankCounts)

	isFlush := false
	flushSuit := Clubs
	for suit, count := range suitCounts {
		if count == 5 {
			isFlush = true
			flushSuit = suit
			break
		}
	}

	if isFlush {
		flushRanks := make([]Rank, 0, 5)
		for _, c := range cards {
			if c.Suit == flushSuit {
				flushRanks = append(flushRanks, c.Rank)
			}
		}
		sort.Slice(flushRanks, func(i, j int) bool { return flushRanks[i] > flushRanks[j] })

		rankMaskFlush := 0
		for _, r := range flushRanks {
			rankMaskFlush |= 1 << int(r)
		}
		if high, ok := straightHighRank(rankMaskFlush); ok {
			return handRankFromSlice(StraightFlush, []Rank{high})
		}

		return handRankFromSlice(Flush, flushRanks)
	}

	if high, ok := straightHighRank(rankMask); ok {
		return handRankFromSlice(Straight, []Rank{high})
	}

	counts := make([]countRank, 0, len(rankCounts))
	for rank, count := range rankCounts {
		counts = append(counts, countRank{Rank: rank, Count: count})
	}
	sort.Slice(counts, func(i, j int) bool {
		if counts[i].Count == counts[j].Count {
			return counts[i].Rank > counts[j].Rank
		}
		return counts[i].Count > counts[j].Count
	})

	if counts[0].Count == 4 {
		fourRank := counts[0].Rank
		kicker := highestExcluding(allRanksDesc, fourRank)
		return handRankFromSlice(FourOfAKind, []Rank{fourRank, kicker})
	}

	if counts[0].Count == 3 {
		tripRank := counts[0].Rank
		if pairRank, ok := highestPairRank(counts[1:]); ok {
			return handRankFromSlice(FullHouse, []Rank{tripRank, pairRank})
		}

		kickers := topRanksExcluding(allRanksDesc, []Rank{tripRank}, 2)
		values := append([]Rank{tripRank}, kickers...)
		return handRankFromSlice(ThreeOfAKind, values)
	}

	pairs := collectPairs(counts)
	if len(pairs) >= 2 {
		first, second := pairs[0], pairs[1]
		if second > first {
			first, second = second, first
		}
		kicker := highestExcluding(allRanksDesc, first, second)
		values := []Rank{first, second, kicker}
		return handRankFromSlice(TwoPair, values)
	}

	if len(pairs) == 1 {
		pair := pairs[0]
		kickers := topRanksExcluding(allRanksDesc, []Rank{pair}, 3)
		values := append([]Rank{pair}, kickers...)
		return handRankFromSlice(OnePair, values)
	}

	highCards := topRanksExcluding(allRanksDesc, nil, 5)
	return handRankFromSlice(HighCard, highCards)
}

type countRank struct {
	Rank  Rank
	Count int
}

func sortedRanksByCount(rankCounts map[Rank]int) []Rank {
	ranks := make([]Rank, 0, len(rankCounts))
	for rank, count := range rankCounts {
		for i := 0; i < count; i++ {
			ranks = append(ranks, rank)
		}
	}
	sort.Slice(ranks, func(i, j int) bool { return ranks[i] > ranks[j] })
	return ranks
}

func highestExcluding(ranks []Rank, excludes ...Rank) Rank {
	excludeSet := make(map[Rank]struct{}, len(excludes))
	for _, e := range excludes {
		excludeSet[e] = struct{}{}
	}
	for _, r := range ranks {
		if _, skip := excludeSet[r]; skip {
			continue
		}
		return r
	}
	return Two
}

func topRanksExcluding(ranks []Rank, excludes []Rank, needed int) []Rank {
	excludeSet := make(map[Rank]struct{}, len(excludes))
	for _, e := range excludes {
		excludeSet[e] = struct{}{}
	}

	result := make([]Rank, 0, needed)
	for _, r := range ranks {
		if _, skip := excludeSet[r]; skip {
			continue
		}
		result = append(result, r)
		if len(result) == needed {
			break
		}
	}
	return result
}

func highestPairRank(counts []countRank) (Rank, bool) {
	for _, cr := range counts {
		if cr.Count >= 2 {
			return cr.Rank, true
		}
	}
	return 0, false
}

func collectPairs(counts []countRank) []Rank {
	pairs := make([]Rank, 0, len(counts))
	for _, cr := range counts {
		if cr.Count >= 2 {
			pairs = append(pairs, cr.Rank)
		}
	}
	return pairs
}

func straightHighRank(mask int) (Rank, bool) {
	for high := Ace; high >= Five; high-- {
		needed := 0
		valid := true
		for offset := 0; offset < 5; offset++ {
			rankValue := int(high) - offset
			if rankValue < 0 {
				valid = false
				break
			}
			needed |= 1 << rankValue
		}
		if valid && mask&needed == needed {
			return high, true
		}
	}

	wheelMask := (1 << int(Ace)) | (1 << int(Two)) | (1 << int(Three)) | (1 << int(Four)) | (1 << int(Five))
	if mask&wheelMask == wheelMask {
		return Five, true
	}

	return 0, false
}

func handRankFromSlice(category HandCategory, ranks []Rank) HandRank {
	h := HandRank{Category: category}
	for i := 0; i < len(ranks) && i < len(h.Values); i++ {
		h.Values[i] = ranks[i]
	}
	return h
}
