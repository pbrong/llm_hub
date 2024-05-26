package llm_caller

import (
	"context"
	"encoding/json"
	"fmt"
	"llm_hub/conf"
	"llm_hub/pkg/http"
)

const (
	// 请求路径
	CompletionsURL = "/v1/chat/completions"

	// 可用模型
	Gpt35TurboModel = "gpt-3.5-turbo"
)

type gptLLMCaller struct {
	openAiKey   string
	systemText  string
	temperature float64
	maxTokens   int64
}

func NewGptLLMCaller(ctx context.Context, systemText string, temperature float64, maxTokens int64) (*gptLLMCaller, error) {
	return &gptLLMCaller{
		openAiKey:   conf.LLMHubConfig.Openai.Key,
		systemText:  systemText,
		temperature: temperature,
		maxTokens:   maxTokens,
	}, nil
}

func (caller *gptLLMCaller) Call(ctx context.Context, userPrompt string) (completion string, err error) {
	reqURL := conf.LLMHubConfig.Openai.Host + CompletionsURL
	body := map[string]interface{}{
		"model":       Gpt35TurboModel,
		"temperature": caller.temperature,
		"stream":      false,
		"max_tokens":  caller.maxTokens,
		"messages":    nil,
	}
	body["messages"] = buildPromptMessages(caller.systemText, userPrompt)
	headers := buildAuthHeaders(caller.openAiKey)
	resp, err := http.PostWithHeader(reqURL, body, headers)
	if err != nil {
		return "", fmt.Errorf("GPT调用失败, err = %v", err)
	}
	var gptCompletion GptCompletion
	_ = json.Unmarshal(resp, &gptCompletion)
	if len(gptCompletion.Choices) > 0 {
		completion = gptCompletion.Choices[0].Message.Content
	}

	return completion, nil
}

func buildAuthHeaders(key string) map[string]string {
	headers := map[string]string{
		"Authorization": "Bearer " + key,
	}
	return headers
}

func buildPromptMessages(system string, user string) []*Message {
	var messages []*Message
	messages = append(messages, &Message{
		Role:    "system",
		Content: system,
	})
	messages = append(messages, &Message{
		Role:    "user",
		Content: user,
	})
	return messages
}
