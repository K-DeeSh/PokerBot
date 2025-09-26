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
Отправьте параметры в формате:
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

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			handleCommand(api, update.Message)
			continue
		}

		text := strings.TrimSpace(update.Message.Text)
		if text == "" {
			continue
		}

		req, err := bot.ParseRequest(text)
		if err != nil {
			reply := tgbotapi.NewMessage(update.Message.Chat.ID, formatError(err))
			reply.ReplyToMessageID = update.Message.MessageID
			sendMessage(api, reply)
			continue
		}

		cfg := req.ToSimulationConfig()
		cfg.Seed = time.Now().UnixNano()

		result, err := poker.SimulateWinProbability(cfg)
		if err != nil {
			reply := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Ошибка симуляции: %v", err))
			reply.ReplyToMessageID = update.Message.MessageID
			sendMessage(api, reply)
			continue
		}

		response := bot.FormatResult(req, result)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
		msg.ReplyToMessageID = update.Message.MessageID
		sendMessage(api, msg)
	}
}

func handleCommand(api *tgbotapi.BotAPI, msg *tgbotapi.Message) {
	reply := tgbotapi.NewMessage(msg.Chat.ID, helpText)
	reply.ReplyToMessageID = msg.MessageID
	sendMessage(api, reply)
}

func formatError(err error) string {
	return fmt.Sprintf("Ошибка: %v\n\n%s", err, helpText)
}

func sendMessage(api *tgbotapi.BotAPI, msg tgbotapi.MessageConfig) {
	if _, err := api.Send(msg); err != nil {
		log.Printf("ошибка отправки сообщения: %v", err)
	}
}
