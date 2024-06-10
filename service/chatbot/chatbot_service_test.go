package chatbot_service

import (
	"context"
	"github.com/stretchr/testify/require"
	"llm_hub/conf"
	"llm_hub/module/chatbot"
	"testing"
)

func Test_chatBotServiceV1_Chat(t *testing.T) {
	conf.Init()
	chatbot.Init()
	ctx := context.Background()
	chatbotService := &chatBotServiceV1{}
	// 模拟用户咨询
	userConsults := []string{
		"你们有哪些产品？",
		"我可以咨询下现在天气怎么样？",
		"电视相关的产品有哪些？",
		"CineView OLED TV和SoundMax Home Theater有什么区别，介绍一下？",
	}
	for idx, consult := range userConsults {
		output, err := chatbotService.Chat(ctx, consult)
		require.Nil(t, err)
		t.Logf("用户咨询%v:%v\n机器人回答:%v", idx+1, consult, output)
	}
}
