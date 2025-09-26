package bot

import (
	"fmt"
	"strconv"
	"strings"

	"pokerbot/internal/poker"
)

// Request captures user intent derived from the incoming message.
type Request struct {
	Hand    []poker.Card
	Board   []poker.Card
	Players int
	Style   poker.PlayerStyle
	Trials  int
}

var styleAliases = map[string]poker.PlayerStyle{
	"balanced":         poker.StyleBalanced,
	"default":          poker.StyleBalanced,
	"neutral":          poker.StyleBalanced,
	"сбалансированный": poker.StyleBalanced,
	"tight":            poker.StyleTight,
	"тайтовый":         poker.StyleTight,
	"тайт":             poker.StyleTight,
	"loose":            poker.StyleLoose,
	"лузовый":          poker.StyleLoose,
	"луз":              poker.StyleLoose,
}

// ParseRequest parses a human-friendly multi-line message into a structured request.
func ParseRequest(text string) (Request, error) {
	lines := strings.Split(text, "\n")
	req := Request{Style: poker.StyleBalanced, Trials: 7000}

	for _, line := range lines {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := normalize(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "hand", "карты":
			hand, err := parseCards(value)
			if err != nil {
				return Request{}, fmt.Errorf("hand: %w", err)
			}
			if len(hand) != 2 {
				return Request{}, fmt.Errorf("hand: expected 2 cards, got %d", len(hand))
			}
			req.Hand = hand
		case "board", "борд", "стол":
			board, err := parseCards(value)
			if err != nil {
				return Request{}, fmt.Errorf("board: %w", err)
			}
			if len(board) > 5 {
				return Request{}, fmt.Errorf("board: expected up to 5 cards, got %d", len(board))
			}
			req.Board = board
		case "players", "игроков", "игроки":
			num, err := parseInt(value)
			if err != nil {
				return Request{}, fmt.Errorf("players: %w", err)
			}
			if num < 2 {
				return Request{}, fmt.Errorf("players: value must be at least 2")
			}
			req.Players = num
		case "style", "стиль":
			style := normalize(value)
			if mapped, ok := styleAliases[style]; ok {
				req.Style = mapped
			} else {
				return Request{}, fmt.Errorf("unknown style: %s", value)
			}
		case "trials", "симуляций":
			num, err := parseInt(value)
			if err != nil {
				return Request{}, fmt.Errorf("trials: %w", err)
			}
			if num < 500 {
				return Request{}, fmt.Errorf("trials: value must be >= 500 for stability")
			}
			req.Trials = num
		}
	}

	if len(req.Hand) != 2 {
		return Request{}, fmt.Errorf("hand: two cards are required")
	}
	if req.Players == 0 {
		return Request{}, fmt.Errorf("players: specify number of players at the table")
	}

	return req, nil
}

func parseCards(value string) ([]poker.Card, error) {
	if strings.TrimSpace(value) == "" {
		return nil, nil
	}

	parts := strings.Fields(value)
	cards := make([]poker.Card, 0, len(parts))
	for _, p := range parts {
		card, err := poker.ParseCard(p)
		if err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}
	return cards, nil
}

func parseInt(value string) (int, error) {
	parts := strings.Fields(value)
	if len(parts) == 0 {
		return 0, fmt.Errorf("missing value")
	}
	return strconv.Atoi(parts[0])
}

func normalize(s string) string {
	return strings.TrimSpace(strings.ToLower(s))
}

// ToSimulationConfig converts a bot request into a simulator configuration.
func (r Request) ToSimulationConfig() poker.SimulationConfig {
	return poker.SimulationConfig{
		Hero:      r.Hand,
		Board:     r.Board,
		Opponents: r.Players - 1,
		Style:     r.Style,
		Trials:    r.Trials,
	}
}
