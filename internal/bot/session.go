package bot

import (
	"fmt"

	"pokerbot/internal/poker"
)

// InputStep indicates which field the bot expects next from the user.
type InputStep int

const (
	StepNone InputStep = iota
	StepHand
	StepPlayers
	StepBoard
	StepTrials
)

// Session keeps track of a user's in-progress request via the menu.
type Session struct {
	Request Request
	Await   InputStep
}

// NewSession returns a session initialised with default values.
func NewSession() Session {
	return Session{
		Request: Request{
			Players: 2,
			Style:   poker.StyleBalanced,
			Trials:  7000,
		},
	}
}

// ApplyValue writes user-provided text into the session based on the awaited step.
func (s *Session) ApplyValue(text string) error {
	switch s.Await {
	case StepHand:
		hand, err := parseCards(text)
		if err != nil {
			return fmt.Errorf("hand: %w", err)
		}
		if len(hand) != 2 {
			return fmt.Errorf("hand: ожидаются две карты")
		}
		s.Request.Hand = hand
	case StepPlayers:
		num, err := parseInt(text)
		if err != nil {
			return fmt.Errorf("players: %w", err)
		}
		if num < 2 {
			return fmt.Errorf("players: минимум два игрока")
		}
		s.Request.Players = num
	case StepBoard:
		board, err := parseCards(text)
		if err != nil {
			return fmt.Errorf("board: %w", err)
		}
		if len(board) > 5 {
			return fmt.Errorf("board: максимум пять карт")
		}
		s.Request.Board = board
	case StepTrials:
		num, err := parseInt(text)
		if err != nil {
			return fmt.Errorf("trials: %w", err)
		}
		if num < 500 {
			return fmt.Errorf("trials: минимум 500")
		}
		s.Request.Trials = num
	default:
		return fmt.Errorf("нет ожидаемого ввода")
	}

	s.Await = StepNone
	return nil
}

// HasRequiredFields reports whether the menu request is ready for simulation.
func (s Session) HasRequiredFields() bool {
	return len(s.Request.Hand) == 2 && s.Request.Players >= 2
}
