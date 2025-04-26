package main

import (
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Помилка з .env! %v", err)
	}

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("Помилка, не знайдено значення ключа з /env!")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Назва бота %v", bot.Self.UserName)

	bot.Debug = true

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		chatID := update.Message.Chat.ID
		if update.Message != nil {
			msg := tgbotapi.NewMessage(chatID, "Твоє повідомлення приняте!")
			bot.Send(msg)
		}
	}
}