package bot

import (
	"testing"

	"pokerbot/internal/poker"
)

func TestParseRequest(t *testing.T) {
	input := "hand: Ah Kh\nplayers: 4\nstyle: tight\nboard: Qh Jd 9c\ntrials: 2000"

	req, err := ParseRequest(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(req.Hand) != 2 || req.Hand[0] != poker.MustParseCard("Ah") {
		t.Fatalf("unexpected hand parsed: %#v", req.Hand)
	}
	if req.Players != 4 {
		t.Fatalf("expected 4 players, got %d", req.Players)
	}
	if req.Style != poker.StyleTight {
		t.Fatalf("expected tight style")
	}
	if req.Trials != 2000 {
		t.Fatalf("expected overridden trials")
	}

	cfg := req.ToSimulationConfig()
	if cfg.Opponents != 3 {
		t.Fatalf("expected 3 opponents, got %d", cfg.Opponents)
	}
}

func TestParseRequestErrors(t *testing.T) {
	_, err := ParseRequest("players: 1\nhand: Ah Kh")
	if err == nil {
		t.Fatal("expected error for invalid player count")
	}

	_, err = ParseRequest("hand: Ah\nplayers: 2")
	if err == nil {
		t.Fatal("expected error for missing second card")
	}

	_, err = ParseRequest("hand: Ah Kh\nplayers: 2\nstyle: ultra")
	if err == nil {
		t.Fatal("expected error for unknown style")
	}
}
