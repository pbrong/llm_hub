package main

import (
	"context"
	"fmt"
	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/callbacks"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/embeddings"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/memory"
	"github.com/tmc/langchaingo/prompts"
	"github.com/tmc/langchaingo/schema"
	"github.com/tmc/langchaingo/tools"
	"github.com/tmc/langchaingo/tools/wikipedia"
	"github.com/tmc/langchaingo/vectorstores"
	"github.com/tmc/langchaingo/vectorstores/redisvector"
	"llm_hub/conf"
	"llm_hub/pkg/json"
	"log"
	"time"
)

func main() {
	conf.Init()
	ctx := context.Background()

	// 调用demo
	//demo(ctx)

	// 指定system及json格式化
	//promptWithRoleJSON(ctx)

	// 提示词调用
	//promptTemplate(ctx)

	// 对话链 + 上下文记忆
	//conversationMemory(ctx)

	// 大模型链
	//llmChains(ctx)

	// 顺序链
	//sequenceChains(ctx)

	// embedding生成
	//embeddingCreate(ctx)

	// rag检索增强生成
	//embeddingRag(ctx)

	// agent使用：数学工具以及搜索
	//agent_math_and_search(ctx)

	// 自定义agent
	agent_diy(ctx)
}

type randomNumberTool struct{}

func (r randomNumberTool) Name() string {
	return "随机数计算工具"
}

func (r randomNumberTool) Description() string {
	return "用于获取随机数"
}

func (r randomNumberTool) Call(ctx context.Context, input string) (string, error) {
	return "1024", nil
}

func agent_diy(ctx context.Context) {
	llm, err := openai.New(
		openai.WithModel("gpt-3.5-turbo"),
		openai.WithBaseURL(conf.LLMHubConfig.Openai.Host),
		openai.WithToken(conf.LLMHubConfig.Openai.Key),
	)
	if err != nil {
		log.Fatal(err)
	}

	agentTools := []tools.Tool{
		randomNumberTool{},
	}
	agent := agents.NewOneShotAgent(llm, agentTools)
	executor := agents.NewExecutor(
		agent,
		agentTools,
		agents.WithCallbacksHandler(callbacks.LogHandler{}),
	)
	result, err := chains.Run(ctx, executor, "告诉我一个随机数")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
}

func agent_math_and_search(ctx context.Context) {
	llm, err := openai.New(
		openai.WithModel("gpt-3.5-turbo"),
		openai.WithBaseURL(conf.LLMHubConfig.Openai.Host),
		openai.WithToken(conf.LLMHubConfig.Openai.Key),
	)
	if err != nil {
		log.Fatal(err)
	}
	wikiTool := wikipedia.New("test")
	agentTools := []tools.Tool{
		tools.Calculator{},
		wikiTool,
	}
	agent := agents.NewOneShotAgent(llm, agentTools)
	executor := agents.NewExecutor(
		agent,
		agentTools,
		agents.WithCallbacksHandler(callbacks.LogHandler{}),
	)
	// 计算
	result, err := chains.Run(ctx, executor, "计算1024除以2并加1024的结果")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
	// 搜索
	result, err = chains.Run(ctx, executor, "今天的日期以及中国在去年今天发生了什么大事")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
}

func embeddingRag(ctx context.Context) {
	// embedding生成测试
	llm, err := openai.New(
		openai.WithEmbeddingModel("text-embedding-ada-002"),
		openai.WithBaseURL(conf.LLMHubConfig.Openai.Host),
		openai.WithToken(conf.LLMHubConfig.Openai.Key),
	)
	if err != nil {
		log.Fatal(err)
	}

	// 创建embedder
	openAiEmbedder, err := embeddings.NewEmbedder(llm)
	if err != nil {
		log.Fatal(err)
	}
	// 基于redis存储向量
	redisStore, err := redisvector.New(ctx,
		redisvector.WithConnectionURL(conf.LLMHubConfig.Redis.Url),
		redisvector.WithIndexName("test_vector_idx", true),
		redisvector.WithEmbedder(openAiEmbedder),
	)
	if err != nil {
		log.Fatalln(err)
	}
	// 插入测试数据
	data := []schema.Document{
		{PageContent: "狸花猫", Metadata: nil},
		{PageContent: "金渐层猫", Metadata: nil},
		{PageContent: "松狮犬", Metadata: nil},
	}

	_, err = redisStore.AddDocuments(ctx, data)
	if err != nil {
		log.Fatalln(err)
	}
	docs, err := redisStore.SimilaritySearch(ctx, "猫", 3,
		vectorstores.WithScoreThreshold(0.5),
	)
	fmt.Println(docs)

	// 将vector检索接入chains中
	result, err := chains.Run(
		ctx,
		chains.NewRetrievalQAFromLLM(
			llm,
			vectorstores.ToRetriever(redisStore, 3, vectorstores.WithScoreThreshold(0.8)),
		),
		"有哪些猫?",
	)
	fmt.Println(result)
}

func embeddingCreate(ctx context.Context) {
	// embedding生成测试
	llm, err := openai.New(
		openai.WithEmbeddingModel("text-embedding-ada-002"),
		openai.WithBaseURL(conf.LLMHubConfig.Openai.Host),
		openai.WithToken(conf.LLMHubConfig.Openai.Key),
	)
	if err != nil {
		log.Fatal(err)
	}
	vectors, err := llm.CreateEmbedding(ctx, []string{"chatgpt-3.5"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(vectors)
}

func llmChains(ctx context.Context) {
	llm, err := openai.New(
		openai.WithModel("gpt-4o"),
		openai.WithBaseURL(conf.LLMHubConfig.Openai.Host),
		openai.WithToken(conf.LLMHubConfig.Openai.Key),
	)
	if err != nil {
		log.Fatal(err)
	}
	// 单个输入
	prompt := prompts.NewPromptTemplate(
		`将"""括起来中文翻译为英文输出
				 输入中文:"""{{.text}}"""
				 输出结果中只需要有英文翻译，不要有其他字符`,
		[]string{"text"})
	llmChain := chains.NewLLMChain(llm, prompt)
	out, err := chains.Run(ctx, llmChain, "langchain是一款不错的llm脚手架")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(out)

	// 多个输入
	translatePrompt := prompts.NewPromptTemplate(
		"Translate the following text from {{.inputLanguage}} to {{.outputLanguage}}. {{.text}}",
		[]string{"inputLanguage", "outputLanguage", "text"},
	)
	llmChain = chains.NewLLMChain(llm, translatePrompt)

	// Otherwise the call function must be used.
	outputValues, err := chains.Call(ctx, llmChain, map[string]any{
		"inputLanguage":  "English",
		"outputLanguage": "Chinese",
		"text":           "I love programming.",
	})
	if err != nil {
		log.Fatal(err)
	}

	out, ok := outputValues[llmChain.OutputKey].(string)
	if !ok {
		log.Fatal(err)
	}
	fmt.Println(out)
}

func routerChains(ctx context.Context) {

}

func sequenceChains(ctx context.Context) {
	llm, err := openai.New(
		openai.WithModel("gpt-4o"),
		openai.WithBaseURL(conf.LLMHubConfig.Openai.Host),
		openai.WithToken(conf.LLMHubConfig.Openai.Key),
	)
	if err != nil {
		log.Fatal(err)
	}

	// 将输入翻译为特定语言
	chain1 := chains.NewLLMChain(llm,
		prompts.NewPromptTemplate(
			"请将输入的原始文本:{{.originText}}翻译为{{.language}}，直接输出翻译文本",
			[]string{"originText", "language"}))
	chain1.OutputKey = "transText"

	// 总结翻译后的文本概要
	chain2 := chains.NewLLMChain(llm, prompts.NewPromptTemplate(
		"请将输入的原始文本:<{{.transText}}>总结50字以内概要文本。严格使用JSON序列化输出结果,不要带有```json序列化标识。其中originText为原始文本,summaryText为概要文本",
		[]string{"transText"}))
	chain2.OutputKey = "summary_json"

	chain, err := chains.NewSequentialChain([]chains.Chain{chain1, chain2}, []string{"originText", "language"}, []string{"summary_json"})
	if err != nil {
		log.Fatal(err)
	}
	resp, err := chain.Call(ctx, map[string]any{
		"originText": "langchain is a good llm frameworks",
		"language":   "中文",
	})
	if err != nil {
		log.Fatal(err)
	}
	for key, value := range resp {
		fmt.Printf("key = %v | value = %v\n", key, value)
	}
}

func conversationMemory(ctx context.Context) {
	llm, err := openai.New(
		openai.WithModel("gpt-4o"),
		openai.WithBaseURL(conf.LLMHubConfig.Openai.Host),
		openai.WithToken(conf.LLMHubConfig.Openai.Key),
	)
	if err != nil {
		log.Fatal(err)
	}
	//memoryBuffer := memory.NewConversationBuffer()
	memoryBuffer := memory.NewConversationWindowBuffer(10)
	//memoryBuffer := memory.NewConversationTokenBuffer(llm, 1024)
	chatChain := chains.NewConversation(llm, memoryBuffer)
	messages := []string{
		"你好，我叫PBR",
		"你知道我叫什么吗？",
		"你可以解决什么问题？",
	}
	for _, message := range messages {
		completion, err := chains.Run(ctx, chatChain, message)
		for {
			if err == nil {
				break
			}
			time.Sleep(30 * time.Second)
			completion, err = chains.Run(ctx, chatChain, message)
		}
		chatMessages, _ := memoryBuffer.ChatHistory.Messages(ctx)
		fmt.Printf("上下文对话历史:%v\n", json.SafeDump(chatMessages))
		fmt.Printf("输入:%v\n输出:%v\n", message, completion)
	}
}

func promptWithRoleJSON(ctx context.Context) {
	llm, err := openai.New(
		openai.WithModel("gpt-4o"),
		openai.WithBaseURL(conf.LLMHubConfig.Openai.Host),
		openai.WithToken(conf.LLMHubConfig.Openai.Key),
	)
	if err != nil {
		log.Fatal(err)
	}

	messages := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, "你是一个英文翻译员，需要将<>括起来的英文翻译为中文，用JSON格式输出：原始文本、翻译文本"),
		llms.TextParts(llms.ChatMessageTypeHuman, "<hello world>"),
	}
	content, err := llm.GenerateContent(ctx, messages, llms.WithJSONMode())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(content.Choices[0].Content)
}

func promptTemplate(ctx context.Context) {
	llm, err := openai.New(
		openai.WithModel("gpt-3.5-turbo"),
		openai.WithBaseURL(conf.LLMHubConfig.Openai.Host),
		openai.WithToken(conf.LLMHubConfig.Openai.Key),
	)
	if err != nil {
		log.Fatal(err)
	}
	prompt := prompts.PromptTemplate{
		Template:       "你是一个文本翻译员,请将```括起来的原始文本转化为{{.lang}}。原始文本```{{.text}}```",
		InputVariables: []string{"text"},
		PartialVariables: map[string]any{
			"lang": "英语",
		},
		TemplateFormat: prompts.TemplateFormatGoTemplate,
	}
	result, err := prompt.Format(map[string]any{
		"text": "我是中国人",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
	result, err = llm.Call(ctx, result)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result)
}

func demo(ctx context.Context, llm *openai.LLM) {
	llm, err := openai.New(
		openai.WithModel("gpt-3.5-turbo"),
		openai.WithBaseURL("https://api.openai-proxy.com/v1"),
		openai.WithToken(conf.LLMHubConfig.Openai.Key),
	)
	if err != nil {
		log.Fatal(err)
	}
	completion, err := llms.GenerateFromSinglePrompt(ctx,
		llm,
		"hello world！",
		llms.WithTemperature(0),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(completion)
}
