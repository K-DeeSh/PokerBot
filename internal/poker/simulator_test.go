package poker

import (
	"math"
	"testing"
)

func TestSimulateWinProbabilitySum(t *testing.T) {
	cfg := SimulationConfig{
		Hero:      []Card{MustParseCard("Ah"), MustParseCard("As")},
		Board:     nil,
		Opponents: 1,
		Style:     StyleBalanced,
		Trials:    2000,
		Seed:      42,
	}

	result, err := SimulateWinProbability(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	sum := result.Win + result.Tie + result.Lose
	if math.Abs(sum-100) > 0.01 {
		t.Fatalf("probabilities must sum to 100, got %.4f", sum)
	}

	if result.Win <= 70 {
		t.Fatalf("expected aces to win frequently, got %.2f%%", result.Win)
	}
}

func TestSimulateWinProbabilityStyles(t *testing.T) {
	hero := []Card{MustParseCard("2c"), MustParseCard("7d")}
	cfgBase := SimulationConfig{
		Hero:      hero,
		Board:     nil,
		Opponents: 1,
		Trials:    3000,
		Seed:      99,
	}

	looseCfg := cfgBase
	looseCfg.Style = StyleLoose
	loose, err := SimulateWinProbability(looseCfg)
	if err != nil {
		t.Fatalf("loose sim error: %v", err)
	}

	balancedCfg := cfgBase
	balancedCfg.Style = StyleBalanced
	balanced, err := SimulateWinProbability(balancedCfg)
	if err != nil {
		t.Fatalf("balanced sim error: %v", err)
	}

	tightCfg := cfgBase
	tightCfg.Style = StyleTight
	tight, err := SimulateWinProbability(tightCfg)
	if err != nil {
		t.Fatalf("tight sim error: %v", err)
	}

	if !(loose.Win >= balanced.Win && balanced.Win >= tight.Win) {
		t.Fatalf("expected win probabilities to order loose >= balanced >= tight, got loose %.2f balanced %.2f tight %.2f", loose.Win, balanced.Win, tight.Win)
	}
}

func TestSimulateWinProbabilityValidation(t *testing.T) {
	_, err := SimulateWinProbability(SimulationConfig{Hero: []Card{MustParseCard("Ah")}})
	if err == nil {
		t.Fatal("expected error for missing hole cards")
	}

	_, err = SimulateWinProbability(SimulationConfig{
		Hero:      []Card{MustParseCard("Ah"), MustParseCard("Kh")},
		Board:     []Card{MustParseCard("2c"), MustParseCard("2d"), MustParseCard("2h"), MustParseCard("2s"), MustParseCard("3c"), MustParseCard("5d")},
		Opponents: 1,
	})
	if err == nil {
		t.Fatal("expected error for oversized board")
	}

	_, err = SimulateWinProbability(SimulationConfig{
		Hero:      []Card{MustParseCard("Ah"), MustParseCard("Ah")},
		Opponents: 1,
	})
	if err == nil {
		t.Fatal("expected error for duplicate cards")
	}
}
