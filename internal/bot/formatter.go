package bot

import (
	"fmt"
	"strings"

	"pokerbot/internal/poker"
)

// FormatResult produces a user-facing reply describing the simulation outcome.
func FormatResult(req Request, result poker.SimulationResult) string {
	var b strings.Builder
	b.WriteString("Вероятности:\n")
	fmt.Fprintf(&b, "Победа: %.2f%%\n", result.Win)
	fmt.Fprintf(&b, "Ничья: %.2f%%\n", result.Tie)
	fmt.Fprintf(&b, "Поражение: %.2f%%\n\n", result.Lose)

	fmt.Fprintf(&b, "Игроков за столом: %d (оппонентов: %d)\n", req.Players, req.Players-1)
	fmt.Fprintf(&b, "Стиль соперников: %s\n", styleDisplay(req.Style))
	fmt.Fprintf(&b, "Симуляций: %d\n", req.Trials)
	fmt.Fprintf(&b, "Ваши карты: %s\n", cardsToText(req.Hand))
	if len(req.Board) > 0 {
		fmt.Fprintf(&b, "Карты на столе: %s\n", cardsToText(req.Board))
	} else {
		b.WriteString("Карты на столе: пока нет\n")
	}

	return b.String()
}

func cardsToText(cards []poker.Card) string {
	parts := make([]string, len(cards))
	for i, c := range cards {
		parts[i] = c.String()
	}
	return strings.Join(parts, " ")
}

func styleDisplay(style poker.PlayerStyle) string {
	switch style {
	case poker.StyleTight:
		return "тайтовый"
	case poker.StyleLoose:
		return "лузовый"
	default:
		return "сбалансированный"
	}
}
