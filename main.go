package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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

	apiKey := os.Getenv("OPENWEATHER_API_KEY")
	if apiKey == "" {
		log.Fatal("Помилка, не знайдено значченя ключа до API!")
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

		if update.Message == nil {
			continue
		}

		if !update.Message.IsCommand() {
			var message string
			city := update.Message.Text

			url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&lang=uk", city, apiKey)
			

			res, err := http.Get(url)
			if err != nil {
				message = "Сталась помилка при запиті для получення данних!"
				msg := tgbotapi.NewMessage(chatID, message)
				bot.Send(msg)
				continue
			}

			defer res.Body.Close()

			boby, err := io.ReadAll(res.Body)
			if err != nil {
				message = "Помилка, при читанні відповіді"
				msg := tgbotapi.NewMessage(chatID, message)
				bot.Send(msg)
				continue
			}

			var result map[string]interface{}
			err = json.Unmarshal(boby, &result)
			if err != nil {
				panic(err)
			}

			cod := result["cod"]

			if codFloat, ok := cod.(float64); ok && codFloat == 200 {
				mainData := result["main"].(map[string]interface{})
				temp := mainData["temp"].(float64) - 273.15
				feelsLike := mainData["feels_like"].(float64) - 273.15
			
				weatherArray := result["weather"].([]interface{})
				weather := weatherArray[0].(map[string]interface{})
				description := weather["description"].(string)
			
				message = fmt.Sprintf("Температура в місті зараз: %.1f°C, відчувається як %.1f°C\nОпис погоди: %s", temp, feelsLike, description)
			} else {
				message = "Місто не знайдено! Перевірте правильність написання."
			}
			
			msg := tgbotapi.NewMessage(chatID, message)
			bot.Send(msg)
		}

		if update.Message.IsCommand() {
			var message string
			command := update.Message.Command()
			userName := update.Message.From.UserName

			switch command {
				case "start":
					message = fmt.Sprintf("Привіт %v я бот прогноз погоди!, пиши своє місто і я тобі скажу яка погода", userName)
					msg := tgbotapi.NewMessage(chatID, message)
					bot.Send(msg)
				default:
					message = "Помилка, не відома команда"
					msg := tgbotapi.NewMessage(chatID, message)
					bot.Send(msg)
			}	
		}
	}
}