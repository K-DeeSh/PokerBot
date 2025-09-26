package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"pokerbot/internal/bot"
	"pokerbot/internal/poker"
)

const helpText = `Привет! Я бот-покерный калькулятор для Техасского Холдэма.
Используйте /menu, чтобы открыть интерактивный конструктор запроса.

Либо отправьте параметры текстом в формате:
hand: Ah Kh
players: 4
style: tight
board: Qh Jh Td
trials: 7000 (необязательно)

Доступные стили: tight, balanced, loose.`

func main() {
	token := strings.TrimSpace(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN не задан")
	}

	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("не удалось создать бота: %v", err)
	}
	api.Debug = false
	log.Printf("Бот авторизован как @%s", api.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := api.GetUpdatesChan(updateConfig)
	sessions := make(map[int64]*bot.Session)

	for update := range updates {
		if update.CallbackQuery != nil {
			handleCallback(api, update.CallbackQuery, sessions)
			continue
		}

		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			handleCommand(api, update.Message, sessions)
			continue
		}

		handleTextMessage(api, update.Message, sessions)
	}
}

func formatError(err error) string {
	return fmt.Sprintf("Ошибка: %v\n\n%s", err, helpText)
}

func sendMessage(api *tgbotapi.BotAPI, msg tgbotapi.MessageConfig) {
	if _, err := api.Send(msg); err != nil {
		log.Printf("ошибка отправки сообщения: %v", err)
	}
}

func handleCommand(api *tgbotapi.BotAPI, msg *tgbotapi.Message, sessions map[int64]*bot.Session) {
	switch msg.Command() {
	case "start":
		sendHelp(api, msg)
		startSession(api, msg.Chat.ID, sessions)
	case "menu":
		startSession(api, msg.Chat.ID, sessions)
	case "cancel":
		delete(sessions, msg.Chat.ID)
		reply := tgbotapi.NewMessage(msg.Chat.ID, "Конструктор сброшен.")
		reply.ReplyToMessageID = msg.MessageID
		sendMessage(api, reply)
	default:
		sendHelp(api, msg)
	}
}

func sendHelp(api *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	reply := tgbotapi.NewMessage(msg.Chat.ID, helpText)
	reply.ReplyToMessageID = msg.MessageID
	sendMessage(api, reply)
}

func startSession(api *tgbotapi.BotAPI, chatID int64, sessions map[int64]*bot.Session) {
	sess := bot.NewSession()
	sessions[chatID] = &sess
	sendMenu(api, chatID, sessions[chatID])
}

func sendMenu(api *tgbotapi.BotAPI, chatID int64, sess *bot.Session) {
	msg := tgbotapi.NewMessage(chatID, bot.SessionSummary(*sess))
	markup := bot.MenuKeyboard()
	msg.ReplyMarkup = markup
	sendMessage(api, msg)
}

func handleTextMessage(api *tgbotapi.BotAPI, msg *tgbotapi.Message, sessions map[int64]*bot.Session) {
	chatID := msg.Chat.ID
	text := strings.TrimSpace(msg.Text)
	if text == "" {
		return
	}

	if sess, ok := sessions[chatID]; ok && sess != nil && sess.Await != bot.StepNone {
		handleAwaitingInput(api, msg, sess)
		sendMenu(api, chatID, sess)
		return
	}

	req, err := bot.ParseRequest(text)
	if err != nil {
		reply := tgbotapi.NewMessage(chatID, formatError(err))
		reply.ReplyToMessageID = msg.MessageID
		sendMessage(api, reply)
		return
	}

	respondWithSimulation(api, msg, req)
}

func handleAwaitingInput(api *tgbotapi.BotAPI, msg *tgbotapi.Message, sess *bot.Session) {
	if err := sess.ApplyValue(msg.Text); err != nil {
		reply := tgbotapi.NewMessage(msg.Chat.ID, err.Error())
		reply.ReplyToMessageID = msg.MessageID
		sendMessage(api, reply)
		promptForStep(api, msg.Chat.ID, sess.Await)
		return
	}

	ack := tgbotapi.NewMessage(msg.Chat.ID, "Принято!")
	ack.ReplyToMessageID = msg.MessageID
	sendMessage(api, ack)
}

func respondWithSimulation(api *tgbotapi.BotAPI, msg *tgbotapi.Message, req bot.Request) {
	cfg := req.ToSimulationConfig()
	cfg.Seed = time.Now().UnixNano()

	result, err := poker.SimulateWinProbability(cfg)
	if err != nil {
		reply := tgbotapi.NewMessage(msg.Chat.ID, fmt.Sprintf("Ошибка симуляции: %v", err))
		reply.ReplyToMessageID = msg.MessageID
		sendMessage(api, reply)
		return
	}

	response := bot.FormatResult(req, result)
	respMessage := tgbotapi.NewMessage(msg.Chat.ID, response)
	respMessage.ReplyToMessageID = msg.MessageID
	sendMessage(api, respMessage)
}

func handleCallback(api *tgbotapi.BotAPI, cb *tgbotapi.CallbackQuery, sessions map[int64]*bot.Session) {
	chatID := cb.Message.Chat.ID
	sess := sessions[chatID]
	if sess == nil {
		s := bot.NewSession()
		sess = &s
		sessions[chatID] = sess
	}

	data := cb.Data

	switch {
	case data == bot.CallbackSetHand:
		sess.Await = bot.StepHand
		promptForStep(api, chatID, bot.StepHand)
	case data == bot.CallbackSetPlayers:
		sess.Await = bot.StepPlayers
		promptForStep(api, chatID, bot.StepPlayers)
	case data == bot.CallbackSetBoard:
		sess.Await = bot.StepBoard
		promptForStep(api, chatID, bot.StepBoard)
	case data == bot.CallbackSetTrials:
		sess.Await = bot.StepTrials
		promptForStep(api, chatID, bot.StepTrials)
	case strings.HasPrefix(data, bot.CallbackSetStyle):
		if style, ok := bot.ParseStyleCallback(data); ok {
			sess.Request.Style = style
			sendMenu(api, chatID, sess)
		} else {
			promptStyleSelection(api, chatID)
		}
	case data == bot.CallbackSetStyle:
		promptStyleSelection(api, chatID)
	case data == bot.CallbackSimulate:
		if !sess.HasRequiredFields() {
			reply := tgbotapi.NewMessage(chatID, "Сначала заполните карты и количество игроков.")
			sendMessage(api, reply)
			break
		}

		req := sess.Request
		req.Trials = sess.Request.Trials
		respondWithSimulation(api, cb.Message, req)
	case data == bot.CallbackCancel:
		delete(sessions, chatID)
		reply := tgbotapi.NewMessage(chatID, "Конструктор очищен. Используйте /menu для нового запроса.")
		sendMessage(api, reply)
	default:
		reply := tgbotapi.NewMessage(chatID, "Неизвестное действие")
		sendMessage(api, reply)
	}

	if _, err := api.Request(tgbotapi.NewCallback(cb.ID, "")); err != nil {
		log.Printf("ошибка ответа на callback: %v", err)
	}
}

func promptForceReply(api *tgbotapi.BotAPI, chatID int64, text, placeholder string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = tgbotapi.ForceReply{ForceReply: true, InputFieldPlaceholder: placeholder}
	sendMessage(api, msg)
}

func promptForStep(api *tgbotapi.BotAPI, chatID int64, step bot.InputStep) {
	var text, placeholder string
	switch step {
	case bot.StepHand:
		text = "Введите две карты героя (например: Ah Kh)"
		placeholder = "Ah Kh"
	case bot.StepPlayers:
		text = "Сколько игроков за столом?"
		placeholder = "4"
	case bot.StepBoard:
		text = "Введите известные карты борда (можно оставить пустым)"
		placeholder = "Qh Jh Th"
	case bot.StepTrials:
		text = "Сколько симуляций выполнить?"
		placeholder = "7000"
	default:
		text = "Введите значение"
		placeholder = ""
	}
	promptForceReply(api, chatID, text, placeholder)
}

func promptStyleSelection(api *tgbotapi.BotAPI, chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Выберите стиль соперников:")
	markup := bot.StyleKeyboard()
	msg.ReplyMarkup = markup
	sendMessage(api, msg)
}
