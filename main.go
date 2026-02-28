package main

import (
	"log"
	"strings"

	"go-chatbot/ai"
	"go-chatbot/config"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	botInfo, err := bot.GetMe()
	if err != nil {
		log.Fatalf("Failed to get bot info: %v", err)
	}

	botUsername := botInfo.UserName
	mention := "@" + botUsername

	log.Printf("Bot @%s started", botUsername)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatalf("Failed to get updates: %v", err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if update.Message.From.ID == botInfo.ID {
			continue
		}

		text := update.Message.Text
		if text == "" {
			continue
		}

		chatType := update.Message.Chat.Type
		isPrivate := chatType == "private"

		// Check for mention (text or entity)
		hasMention := strings.Contains(text, mention)
		
		if !hasMention && update.Message.Entities != nil {
			for _, entity := range *update.Message.Entities {
				if entity.Type == "mention" || entity.Type == "text_mention" {
					hasMention = true
					break
				}
			}
		}

		// Check if replying to bot's message
		if !hasMention && update.Message.ReplyToMessage != nil {
			if update.Message.ReplyToMessage.From.ID == botInfo.ID {
				hasMention = true
			}
		}

		// Skip group messages without mention
		if !isPrivate && !hasMention {
			continue
		}

		log.Printf("[%s] %s", chatType, text)

		// Clean prompt
		prompt := strings.ReplaceAll(text, mention, "")
		prompt = strings.TrimSpace(strings.Trim(prompt, ",.!? "))

		if prompt == "" {
			prompt = "Hello!"
		}

		// Send typing indicator
		typing := tgbotapi.NewChatAction(update.Message.Chat.ID, tgbotapi.ChatTyping)
		bot.Send(typing)

		// Get AI response
		response, err := ai.Chat(cfg.OllamaModel, cfg.SystemPrompt, prompt)
		if err != nil {
			log.Printf("[AI ERROR] %v", err)
			response = "Sorry, I'm having trouble responding."
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
		msg.ReplyToMessageID = update.Message.MessageID

		if _, err := bot.Send(msg); err != nil {
			log.Printf("[SEND ERROR] %v", err)
		} else {
			log.Printf("[SENT] %d chars", len(response))
		}
	}
}
