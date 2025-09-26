package poker

import (
	"errors"
	"math/rand"
	"time"
)

// PlayerStyle categorises opponent tendencies used in the simulation.
type PlayerStyle int

const (
	StyleBalanced PlayerStyle = iota
	StyleTight
	StyleLoose
)

// SimulationConfig describes the parameters for a Monte Carlo probability calculation.
type SimulationConfig struct {
	Hero      []Card
	Board     []Card
	Opponents int
	Style     PlayerStyle
	Trials    int
	Seed      int64
}

// SimulationResult contains aggregate probabilities.
type SimulationResult struct {
	Win  float64
	Tie  float64
	Lose float64
}

// SimulateWinProbability estimates hero equity via Monte Carlo sampling.
func SimulateWinProbability(cfg SimulationConfig) (SimulationResult, error) {
	if len(cfg.Hero) != 2 {
		return SimulationResult{}, errors.New("hero must have exactly two hole cards")
	}
	if len(cfg.Board) > 5 {
		return SimulationResult{}, errors.New("board cannot exceed five cards")
	}
	if cfg.Opponents < 1 {
		return SimulationResult{}, errors.New("opponent count must be at least one")
	}
	if cfg.Opponents > 8 {
		return SimulationResult{}, errors.New("opponent count too large (max 8)")
	}

	distinct := make(map[Card]struct{}, len(cfg.Hero)+len(cfg.Board))
	for _, c := range append(append([]Card(nil), cfg.Hero...), cfg.Board...) {
		if _, exists := distinct[c]; exists {
			return SimulationResult{}, errors.New("duplicate cards provided")
		}
		distinct[c] = struct{}{}
	}

	trials := cfg.Trials
	if trials <= 0 {
		trials = 5000
	}

	seed := cfg.Seed
	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	rng := rand.New(rand.NewSource(seed))

	wins, ties, losses := 0, 0, 0
	excluded := append([]Card(nil), cfg.Hero...)
	excluded = append(excluded, cfg.Board...)

	for i := 0; i < trials; i++ {
		deck := BuildDeck(excluded)
		rng.Shuffle(len(deck), func(i, j int) { deck[i], deck[j] = deck[j], deck[i] })

		board := make([]Card, len(cfg.Board))
		copy(board, cfg.Board)
		if len(board) < 5 {
			needed := 5 - len(board)
			board = append(board, deck[:needed]...)
			deck = deck[needed:]
		}

		heroCards := append(append([]Card(nil), cfg.Hero...), board...)
		heroRank, err := EvaluateBestHand(heroCards)
		if err != nil {
			return SimulationResult{}, err
		}

		heroWins := true
		heroTies := false

		for opp := 0; opp < cfg.Opponents; opp++ {
			hand := drawOpponentHand(&deck, cfg.Style, rng)
			if len(hand) != 2 {
				return SimulationResult{}, errors.New("not enough cards to draw opponent hand")
			}

			oppCards := append(append([]Card(nil), hand...), board...)
			oppRank, err := EvaluateBestHand(oppCards)
			if err != nil {
				return SimulationResult{}, err
			}

			switch heroRank.Compare(oppRank) {
			case -1:
				heroWins = false
				heroTies = false
				goto outcome
			case 0:
				heroTies = true
			}
		}

	outcome:
		if !heroWins {
			losses++
		} else if heroTies {
			ties++
		} else {
			wins++
		}
	}

	return SimulationResult{
		Win:  percentage(wins, trials),
		Tie:  percentage(ties, trials),
		Lose: percentage(losses, trials),
	}, nil
}

func drawOpponentHand(deck *[]Card, style PlayerStyle, rng *rand.Rand) []Card {
	cards := *deck
	if len(cards) < 2 {
		return nil
	}

	window := len(cards)
	if window > 6 {
		window = 6
	}
	bestI, bestJ := 0, 1
	bestScore := startingHandScore(cards[0], cards[1])

	switch style {
	case StyleBalanced:
		// keep defaults
	case StyleTight:
		for i := 0; i < window-1; i++ {
			for j := i + 1; j < window; j++ {
				score := startingHandScore(cards[i], cards[j])
				if score > bestScore || (score == bestScore && rng.Intn(2) == 0) {
					bestScore = score
					bestI, bestJ = i, j
				}
			}
		}
	case StyleLoose:
		for i := 0; i < window-1; i++ {
			for j := i + 1; j < window; j++ {
				score := startingHandScore(cards[i], cards[j])
				if score < bestScore || (score == bestScore && rng.Intn(2) == 0) {
					bestScore = score
					bestI, bestJ = i, j
				}
			}
		}
	}

	var selected [2]Card
	selected[0] = cards[bestI]
	selected[1] = cards[bestJ]

	if bestI > bestJ {
		bestI, bestJ = bestJ, bestI
	}

	cards = append(cards[:bestJ], cards[bestJ+1:]...)
	cards = append(cards[:bestI], cards[bestI+1:]...)

	*deck = cards
	return selected[:]
}

func startingHandScore(a, b Card) int {
	score := 0
	if a.Rank == b.Rank {
		score += 400 + int(a.Rank)*20
	}

	score += (int(a.Rank) + int(b.Rank)) * 10

	gap := int(a.Rank) - int(b.Rank)
	if gap < 0 {
		gap = -gap
	}
	score -= gap * 4

	if a.Suit == b.Suit {
		score += 30
	}

	return score
}

func percentage(count, total int) float64 {
	if total == 0 {
		return 0
	}
	return float64(count) * 100 / float64(total)
}
