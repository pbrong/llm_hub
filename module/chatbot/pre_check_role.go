package chatbot

import (
	"context"
	"fmt"
	"llm_hub/pkg/llm_caller"
)

type preCheckRole struct{}

func NewPreCheckRole() ChatRole {
	return &preCheckRole{}
}

func (role *preCheckRole) Chat(ctx context.Context, input string, optArg *OptArg) (output string, err error) {
	systemPrompt := `
		请仔细检查{}括起来的用户输入文本，判断该输入文本是否与电商咨询/售前咨询问题相关。
		你只可以输出Y或N，分别代表是相关咨询问题和不是相关咨询问题。
`
	gptCaller, err := llm_caller.NewGptLLMCaller(ctx, systemPrompt, 0, 1024)
	if err != nil {
		return "", err
	}
	completion, err := gptCaller.Call(ctx, fmt.Sprintf("用户输入:{%v}", input))
	if err != nil {
		return "", err
	}

	return completion, nil
}
