package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Request struct {
	Model    string `json:"model"`
	Messages []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"messages"`
	Stream bool `json:"stream"`
}

type Response struct {
	Message struct {
		Content string `json:"content"`
	} `json:"message"`
	Done bool `json:"done"`
}

func Chat(model, systemPrompt, prompt string) (string, error) {
	messages := []struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: prompt},
	}

	req := Request{
		Model:    model,
		Messages: messages,
		Stream:   false,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	resp, err := http.Post("http://localhost:11434/api/chat", "application/json", bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("Ollama error: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status %d: %s", resp.StatusCode, body)
	}

	var result Response
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	return strings.TrimSpace(result.Message.Content), nil
}
