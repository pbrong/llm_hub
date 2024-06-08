package chatbot

import (
	"context"
	"fmt"
	"llm_hub/pkg/llm_caller"
)

type outputOptRole struct{}

func NewOutputOptRole() ChatRole {
	return &outputOptRole{}
}

func (role *outputOptRole) Chat(ctx context.Context, input string, optArg *OptArg) (output string, err error) {
	systemPrompt := `
		你是一个友好的客服机器人，现在需要对{}括起来的原始回答信息进行优化，并输出优化后的回答。
		要求：
		1-优化后的回答不能与原始回答存在歧义
		2-优化后的回答需要显示出耐心、友好和语句通顺
		3-将原始回答中的JSON格式优化为可以阅读的自然语言格式输出
		4-输出的优化回答中，不要带有经过优化等字样
`
	gptCaller, err := llm_caller.NewGptLLMCaller(ctx, systemPrompt, 0, 2048)
	if err != nil {
		return "", err
	}
	completion, err := gptCaller.Call(ctx, fmt.Sprintf("原始回答信息:{%v}", input))
	if err != nil {
		return "", err
	}

	return completion, nil
}
