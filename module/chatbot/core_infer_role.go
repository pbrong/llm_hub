package chatbot

import (
	"context"
	"errors"
	"fmt"
	"llm_hub/pkg/llm_caller"
)

type coreInferRole struct{}

func NewCoreInferRole() ChatRole {
	return &coreInferRole{}
}

func (role *coreInferRole) Chat(ctx context.Context, input string, optArg *OptArg) (output string, err error) {
	if len(optArg.Products) == 0 {
		return "", errors.New("product not exist")
	}
	systemPrompt := `
		请仔细阅读并理解{}括起来的用户咨询问题，给出你的咨询回答。请严格遵循以下过程判断和回答：
		1-首先判断用户咨询的问题中，是否存在商品相关信息
		2-若存在商品相关信息，则结合{}括起来的商品信息进行回答
`
	gptCaller, err := llm_caller.NewGptLLMCaller(ctx, systemPrompt, 0, 2048)
	if err != nil {
		return "", err
	}
	completion, err := gptCaller.Call(ctx, fmt.Sprintf("用户咨询问题:{%v}\n商品信息:{%v}", input, optArg.Products))
	if err != nil {
		return "", err
	}

	return completion, nil
}
