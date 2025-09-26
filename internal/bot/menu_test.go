package bot

import (
	"strings"
	"testing"

	"pokerbot/internal/poker"
)

func TestParseStyleCallback(t *testing.T) {
	cases := []struct {
		data  string
		valid bool
		style poker.PlayerStyle
	}{
		{styleCallback(pokerStyleBalanced), true, poker.StyleBalanced},
		{styleCallback(pokerStyleTight), true, poker.StyleTight},
		{styleCallback(pokerStyleLoose), true, poker.StyleLoose},
		{"other", false, poker.StyleBalanced},
	}

	for _, tc := range cases {
		style, ok := ParseStyleCallback(tc.data)
		if ok != tc.valid || (ok && style != tc.style) {
			t.Fatalf("unexpected result for %s", tc.data)
		}
	}
}

func TestSessionSummary(t *testing.T) {
	sess := NewSession()
	sess.Request.Hand = []poker.Card{poker.MustParseCard("Ah"), poker.MustParseCard("Kh")}
	sess.Request.Players = 4
	sess.Request.Board = []poker.Card{poker.MustParseCard("Qh"), poker.MustParseCard("Jh"), poker.MustParseCard("Th")}

	summary := SessionSummary(sess)
	for _, fragment := range []string{"Ah Kh", "Игроки: 4", "Борд: Qh Jh Th", "Запустить"} {
		if !strings.Contains(summary, fragment) {
			t.Fatalf("expected summary to contain %q", fragment)
		}
	}
}
