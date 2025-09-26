package bot

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"pokerbot/internal/poker"
)

const (
	CallbackSetHand    = "set_hand"
	CallbackSetPlayers = "set_players"
	CallbackSetBoard   = "set_board"
	CallbackSetTrials  = "set_trials"
	CallbackSetStyle   = "set_style"
	CallbackSimulate   = "simulate"
	CallbackCancel     = "cancel"
)

// MenuKeyboard returns inline keyboard markup for the interactive builder.
func MenuKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Карты", CallbackSetHand),
			tgbotapi.NewInlineKeyboardButtonData("Игроки", CallbackSetPlayers),
			tgbotapi.NewInlineKeyboardButtonData("Стиль", CallbackSetStyle),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Борд", CallbackSetBoard),
			tgbotapi.NewInlineKeyboardButtonData("Симуляции", CallbackSetTrials),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Запустить", CallbackSimulate),
			tgbotapi.NewInlineKeyboardButtonData("Отмена", CallbackCancel),
		),
	)
}

// SessionSummary renders the current session values for the user.
func SessionSummary(s Session) string {
	var b strings.Builder
	b.WriteString("Конструктор запроса\n")
	b.WriteString("Выберите параметры кнопками ниже.\n\n")
	b.WriteString(formatSessionLine("Карты", cardsDisplay(s.Request.Hand)))
	b.WriteString(formatSessionLine("Игроки", playersDisplay(s.Request.Players)))
	b.WriteString(formatSessionLine("Стиль", styleDisplay(s.Request.Style)))
	b.WriteString(formatSessionLine("Борд", cardsDisplay(s.Request.Board)))
	b.WriteString(formatSessionLine("Симуляций", trialsDisplay(s.Request.Trials)))
	b.WriteString("\nНажмите \"Запустить\", чтобы рассчитать вероятность.")
	return b.String()
}

func formatSessionLine(label, value string) string {
	return fmt.Sprintf("%s: %s\n", label, value)
}

func cardsDisplay(cards []poker.Card) string {
	if len(cards) == 0 {
		return "не задано"
	}
	return CardsToText(cards)
}

func playersDisplay(players int) string {
	if players == 0 {
		return "не задано"
	}
	return fmt.Sprintf("%d", players)
}

func trialsDisplay(trials int) string {
	if trials == 0 {
		return "по умолчанию (7000)"
	}
	return fmt.Sprintf("%d", trials)
}

// StyleKeyboard enumerates style options.
func StyleKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Сбалансированный", styleCallback(pokerStyleBalanced)),
			tgbotapi.NewInlineKeyboardButtonData("Тайтовый", styleCallback(pokerStyleTight)),
			tgbotapi.NewInlineKeyboardButtonData("Лузовый", styleCallback(pokerStyleLoose)),
		),
	)
}

const (
	pokerStyleBalanced = "balanced"
	pokerStyleTight    = "tight"
	pokerStyleLoose    = "loose"
)

func styleCallback(val string) string {
	return CallbackSetStyle + ":" + val
}

// ParseStyleCallback maps callback data to player styles.
func ParseStyleCallback(data string) (poker.PlayerStyle, bool) {
	switch data {
	case styleCallback(pokerStyleBalanced):
		return poker.StyleBalanced, true
	case styleCallback(pokerStyleTight):
		return poker.StyleTight, true
	case styleCallback(pokerStyleLoose):
		return poker.StyleLoose, true
	default:
		return poker.StyleBalanced, false
	}
}
