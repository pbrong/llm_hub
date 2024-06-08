package chatbot

import "context"

var (
	// 前置校验角色
	PreCheckRole ChatRole
	// 核心对话推理角色
	CoreInferRole ChatRole
	// 对话输出优化角色
	OutputOptRole ChatRole
)

type ChatRole interface {
	Chat(ctx context.Context, input string, optArgs *OptArg) (output string, err error)
}

type OptArg struct {
	Products []string
}

type ChatChain []ChatRole

func NewChatChain(roles ...ChatRole) ChatChain {
	chain := ChatChain{}
	chain = append(chain, roles...)
	return chain
}

func (chain ChatChain) Run(ctx context.Context, input string, optArg *OptArg) (output string, err error) {
	for _, role := range chain {
		output, err = role.Chat(ctx, input, optArg)
		if err != nil {
			return "", err
		}
		if output == "N" {
			return "请咨询相关问题", nil
		}
		if output != "" && output != "Y" {
			input = output
		}
	}

	return output, nil
}

func Init() {
	PreCheckRole = NewPreCheckRole()
	CoreInferRole = NewCoreInferRole()
	OutputOptRole = NewOutputOptRole()
}
