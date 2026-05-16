package main

import (
	"fmt"
	"log"
	"strings"

	"go-chatbot/ai"
	"go-chatbot/config"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type ChatMessage struct {
	Role    string
	Content string
}

var history = make(map[int64][]ChatMessage)
var contextLimit int

func addToHistory(userID int64, role, content string) {
	history[userID] = append(history[userID], ChatMessage{Role: role, Content: content})
	if len(history[userID]) > contextLimit*2 {
		history[userID] = history[userID][2:]
	}
}

func buildContextPrompt(userID int64) string {
	if len(history[userID]) == 0 {
		return ""
	}
	var sb strings.Builder
	sb.WriteString("Previous conversation:\n")
	for _, msg := range history[userID] {
		sb.WriteString(fmt.Sprintf("%s: %s\n", msg.Role, msg.Content))
	}
	sb.WriteString("\n")
	return sb.String()
}

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

	contextLimit = cfg.ContextLimit

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

		userID := update.Message.From.ID
		chatID := update.Message.Chat.ID

		log.Printf("[%s] %s", chatType, text)

		var prompt strings.Builder
		prompt.WriteString(buildContextPrompt(int64(userID)))

		userText := strings.ReplaceAll(text, mention, "")
		userText = strings.TrimSpace(strings.Trim(userText, ",.!? "))

		if update.Message.ReplyToMessage != nil && update.Message.ReplyToMessage.Text != "" {
			prompt.WriteString("Replying to: " + update.Message.ReplyToMessage.Text + "\n\n")
		}

		if userText != "" {
			prompt.WriteString(userText)
		} else {
			prompt.WriteString("Hello!")
		}

		// Send typing indicator
		typing := tgbotapi.NewChatAction(chatID, tgbotapi.ChatTyping)
		bot.Send(typing)

		// Get AI response
		response, err := ai.Chat(cfg.OllamaModel, cfg.SystemPrompt, prompt.String())
		if err != nil {
			log.Printf("[AI ERROR] %v", err)
			response = "Sorry, I'm having trouble responding."
		}

		addToHistory(int64(userID), "user", userText)
		addToHistory(int64(userID), "assistant", response)

		msg := tgbotapi.NewMessage(chatID, response)
		msg.ReplyToMessageID = update.Message.MessageID

		if _, err := bot.Send(msg); err != nil {
			log.Printf("[SEND ERROR] %v", err)
		} else {
			log.Printf("[SENT] %d chars", len(response))
		}
	}
}
