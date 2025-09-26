package bot

import (
	"strings"
	"testing"

	"pokerbot/internal/poker"
)

func TestFormatResult(t *testing.T) {
	req := Request{
		Hand:    []poker.Card{poker.MustParseCard("Ah"), poker.MustParseCard("Kh")},
		Board:   []poker.Card{poker.MustParseCard("Qh"), poker.MustParseCard("Jh"), poker.MustParseCard("Th")},
		Players: 4,
		Style:   poker.StyleTight,
		Trials:  5000,
	}

	res := poker.SimulationResult{Win: 55.5, Tie: 3.3, Lose: 41.2}
	text := FormatResult(req, res)

	for _, fragment := range []string{"55.50", "Игроков за столом: 4", "тайтовый", "Ah Kh", "Карты на столе"} {
		if !strings.Contains(text, fragment) {
			t.Fatalf("expected output to contain %q, got: %s", fragment, text)
		}
	}
}
