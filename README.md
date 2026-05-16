# Go Chatbot

A minimal Telegram chatbot powered by Ollama for local AI responses.

## Overview

This bot connects Telegram to your local Ollama instance. It responds to messages in private chats and group chats (when mentioned or replied to).

## Requirements

- Go 1.21 or higher
- Ollama installed and running locally
- Telegram bot token from @BotFather

## Installation

1. Clone or download this repository

2. Install dependencies:
   ```
   go mod tidy
   ```

3. Build:
   ```
   go build -o go-chatbot.exe
   ```

## Configuration

Create a `config.json` file in the same directory as the executable:

```json
{
  "telegram_token": "your_bot_token_here",
  "ollama_model": "llama3.2",
  "system_prompt": "You are a helpful Telegram chatbot. Keep responses short and conversational.",
  "context_limit": 5
}
```

### Configuration Options

| Field | Description | Required | Default |
|-------|-------------|----------|---------|
| `telegram_token` | Bot token from @BotFather | Yes | - |
| `ollama_model` | Ollama model name | No | llama3.2 |
| `system_prompt` | Instructions for AI behavior | No | Default helpful assistant |
| `context_limit` | Max conversation exchanges to remember per user | No | 5 |

## Usage

1. Start Ollama and pull a model:
   ```
   ollama pull llama3.2
   ```

2. Run the bot:
   ```
   ./go-chatbot.exe
   ```

3. In Telegram:
   - **Private chats**: Send any message, the bot replies
   - **Group chats**: Mention `@botname` or reply to the bot's message
   - **Conversation memory**: Bot remembers your last N exchanges across all chats

### Group Privacy Note

By default, Telegram bots cannot see all group messages. To enable full functionality:

1. Message @BotFather
2. Send `/mybots` and select your bot
3. Go to Bot Settings > Group Privacy
4. Turn it OFF

**Note:** In supergroups, the bot must be an admin to receive @mentions reliably.

## System Prompt Examples

Customize how your bot behaves:

```json
"system_prompt": "You are a sarcastic, witty assistant. Make jokes but stay helpful."
```

```json
"system_prompt": "You are a professional customer support agent. Be polite and concise."
```

```json
"system_prompt": "You are a terse assistant. Answer in one sentence only."
```

```json
"system_prompt": "You are an expert programmer. Help with coding questions."
```

## Project Structure

```
go-chatbot/
├── main.go              # Bot logic and message handling
├── config.json          # Your configuration
├── go.mod               # Go module file
├── go.sum               # Dependency checksums
├── config/
│   └── config.go        # Configuration loader
└── ai/
    └── ai.go            # Ollama API client
```

## Logging

The bot outputs to console:

```
Bot @yourbot started
[supergroup] @botname hello
[SENT] 45 chars
[private] test
[SENT] 32 chars
[AI ERROR] Ollama error: ...
```

## Features

- Shows typing indicator while processing
- Responds to mentions in groups
- Responds to replies to bot messages in groups
- Responds to all messages in private chats
- Skips messages from other bots
- Conversation history: remembers recent exchanges per user
- Reply context: includes the message you're replying to

## Troubleshooting

**Bot doesn't respond in groups:**
- Disable Group Privacy in @BotFather, or add bot as admin
- Make sure to mention the bot or reply to its message

**AI errors:**
- Ensure Ollama is running: `ollama serve`
- Check model is pulled: `ollama pull llama3.2`
- Verify model name in config.json

**Bot doesn't start:**
- Check telegram_token is valid
- Ensure config.json is in the same directory as the executable

## License

MIT
