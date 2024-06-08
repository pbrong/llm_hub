package chatbot_service

import (
	"context"
	"llm_hub/module/chatbot"
)

var ChatBotServiceV1 = chatBotServiceV1{}

type ChatBotService interface {
	Chat(ctx context.Context, message string) (output string, err error)
}

type chatBotServiceV1 struct{}

func (service *chatBotServiceV1) Chat(ctx context.Context, message string) (output string, err error) {
	// 初始化角色链
	chain := chatbot.NewChatChain(
		chatbot.PreCheckRole,
		chatbot.CoreInferRole,
		chatbot.OutputOptRole)

	// 检索商品信息(RAG)
	products := retrievalProducts(ctx)

	// 调用llm
	output, err = chain.Run(ctx, message, &chatbot.OptArg{
		Products: products,
	})
	if err != nil {
		return "", err
	}

	return output, nil
}

func retrievalProducts(ctx context.Context) []string {
	// Mock一批数据，这里其实是一段RAG搜索逻辑（by向量匹配）
	return []string{
		`{
			"名称": "SoundMax Soundbar",
			"类别": "电视和家庭影院系统",
			"品牌": "SoundMax",
			"型号": "SM-SB50",
			"保修期": "1 year",
			"评分": 4.3,
			"特色": [
			"2.1 channel",
			"300W output",
			"Wireless subwoofer",
			"Bluetooth"
			],
			"描述": "使用这款时尚而功能强大的声音，升级您电视的音频体验。",
			"价格": 199.99
			}`,
		`{
			"名称": "SoundMax Home Theater",
			"类别": "电视和家庭影院系统",
			"品牌": "SoundMax",
			"型号": "SM-HT100",
			"保修期": "1 year",
			"评分": 4.4,
			"特色": [3.3 生成用户查询的答案
			"5.1 channel",
			"1000W output",
			"Wireless subwoofer",
			"Bluetooth"
			],
			"描述": "一款强大的家庭影院系统，提供沉浸式音频体验。",
			"价格": 399.99
		}`,
		`{
			"名称": "CineView OLED TV",
			"类别": "电视和家庭影院系统",
			"品牌": "CineView",
			"型号": "CV-OLED55",
			"保修期": "2 years",
			"评分": 4.7,
			"特色": [
			"55-inch display",
			"4K resolution",
			"HDR",
			"Smart TV"
			],
			"描述": "通过这款OLED电视，体验真正的五彩斑斓。",
			"价格": 1499.99
		}`,
		`{
			"名称": "CineView 8K TV",
			"类别": "电视和家庭影院系统",
			"品牌": "CineView",
			"型号": "CV-8K65",
			"保修期": "2 years",
			"评分": 4.9,
			"特色": [
			"65-inch display",
			"8K resolution",
			"HDR",
			"Smart TV"
			],
			"描述": "通过这款惊艳的8K电视，体验未来。",
			"价格": 2999.99
		}`,
		`{
			"名称": "CineView 4K TV",
			"类别": "电视和家庭影院系统",
			"品牌": "CineView",
			"型号": "CV-4K55",
			"保修期": "2 years",
			"评分": 4.8,
			"特色": [
			"55-inch display",
			"4K resolution",
			"HDR",
			"Smart TV"
			],
			"描述": "一款色彩鲜艳、智能功能丰富的惊艳4K电视。",
			"价格": 599.99
		}`,
	}
}
