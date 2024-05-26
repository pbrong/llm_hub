package llm_caller

import "context"

var (
	_ LLMCaller = &gptLLMCaller{}
)

type LLMCaller interface {
	Call(ctx context.Context, userPrompt string) (completions string, err error)
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GptCompletion struct {
	Created int `json:"created"`
	Usage   struct {
		CompletionTokens int `json:"completion_tokens"`
		PromptTokens     int `json:"prompt_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Model   string `json:"model"`
	ID      string `json:"id"`
	Choices []struct {
		FinishReason string `json:"finish_reason"`
		Index        int    `json:"index"`
		Message      struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	SystemFingerprint interface{} `json:"system_fingerprint"`
	Object            string      `json:"object"`
}
