package main

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"llm_hub/conf"
	"llm_hub/pkg/llm_caller"
)

func main() {
	conf.Init()
	ctx := context.Background()
	// 使用分隔符
	//promptDelimiter(ctx)

	// json格式化
	//promptJSON(ctx)

	// 检查是否符合条件
	//promptCheck(ctx)

	// few-shot，少样本提示
	//promptFewShot(ctx)

	// 指定完成任务的步骤及结构
	//promptStructure(ctx)

	// 给出结论前先思考并校验，给出推理过程
	promptCoT(ctx)
}

func promptCoT(ctx context.Context) {
	systemText := `你是一个数学老师，请你批改{}括起来的含有作业题目及答案的文本，执行以下操作：
1.首先阅读题目，并自己计算得到答案，写下推理过程
2.将你的推理过程与文本中的答案比对
3.判断文本中的答案是否正确
4.总结文本中的答案错在了哪里

最后，上述操作无需输出，只需要使用JSON格式输出以下Key及值：推理过程、答案是否正确、答案错在哪里

对于JSON格式中Key的说明：
1.推理过程：你解答这个题目的推理过程
2.答案是否正确：输入文本中的答案是否正确，输出值为:正确、不正确	
3.答案错在哪里：输入文本中的答案错在哪里`
	gptCaller, err := llm_caller.NewGptLLMCaller(ctx, systemText, 0.0, 1024)
	if err != nil {
		logrus.Errorf("gpt caller init fail, err = %v,", err)
		return
	}
	text := `题目:
土地费用为 100美元/平方英尺
我可以以 250美元/平方英尺的价格购买太阳能电池板
我已经谈判好了维护合同，每年需要支付固定的10万美元，并额外支付每平方英尺10美元
作为平方英尺数的函数，首年运营的总费用是多少。
答案：
设x为发电站的大小，单位为平方英尺。
费用：
土地费用：100x
太阳能电池板费用：250x
维护费用：100,000美元+100x
总费用：100x+250x+100,000美元+100x=450x+100,000美元`
	resp, err := gptCaller.Call(ctx, fmt.Sprintf("输入文本：{%v}", text))
	if err != nil {
		logrus.Errorf("gpt caller call fail, err = %v", err)
		return
	}
	logrus.Info(resp)
}

func promptStructure(ctx context.Context) {
	systemText := `请你对{}括起来的文本执行以下操作：
1.首先判断该文本属于哪类语言
2.将该文本翻译为英文
3.总结出中文文本中有哪些地点
4.总结出中文文本中有哪些人物

使用JSON形式表达上述操作结果，输出以下Key及值：原始语言、英文翻译、地点汇总、人物汇总`
	gptCaller, err := llm_caller.NewGptLLMCaller(ctx, systemText, 0.0, 512)
	if err != nil {
		logrus.Errorf("gpt caller init fail, err = %v,", err)
		return
	}
	text := `小村庄里，老师小林找到了一把古老的钥匙，决定带领李雷和韩梅梅探索传说中的秘密花园。
他们穿过幽深的森林，解开了一系列谜题，最终发现花园中隐藏的是一座知识的宝库，让整个村庄的孩子们受益无穷。`
	resp, err := gptCaller.Call(ctx, fmt.Sprintf("输入文本：{%v}", text))
	if err != nil {
		logrus.Errorf("gpt caller call fail, err = %v", err)
		return
	}
	logrus.Info(resp)
}

func promptFewShot(ctx context.Context) {
	systemText := `你是一个小红书文案生成器，你需要将{}包裹起来的输入内容改写为小红书的文案风格并输出改写后的文案。以下是几段小红书风格的博文参考：
1、嘴上说不要，钱包却很诚实嘛。
2、原谅我一生放浪不羁爱shopping！
3、草！我又拔草了！
4、工资已经到账，快递还会远吗？
5、掐指一算，姑娘你命中缺钱。
注意，只需要给出一句话作为文案，不要有多余的特殊符号输出`
	gptCaller, err := llm_caller.NewGptLLMCaller(ctx, systemText, 0.0, 512)
	if err != nil {
		logrus.Errorf("gpt caller init fail, err = %v,", err)
		return
	}
	text := `大模型是当今的潮流，学习及应用大模型是一个大趋势，让我们一起学习大模型`
	resp, err := gptCaller.Call(ctx, fmt.Sprintf("输入内容为{%v}", text))
	if err != nil {
		logrus.Errorf("gpt caller call fail, err = %v", err)
		return
	}
	logrus.Info(resp)
}

func promptCheck(ctx context.Context) {
	systemText := `你是一个操作步骤总结器，请阅读{}括起来的输入内容，并提炼出输入内容中提到的执行步骤。
提炼出的步骤按以下格式返回：
- 第一步：...
- 第二步：...
注意，对于不存在明确步骤的输入内容，直接说"不存在步骤"`
	gptCaller, err := llm_caller.NewGptLLMCaller(ctx, systemText, 0.0, 512)
	if err != nil {
		logrus.Errorf("gpt caller init fail, err = %v,", err)
		return
	}
	{
		text := `制作牛排首先需要准备材料有牛肉1块，洋葱半个，蒜头2个，黄油1小块，黑胡椒粉少许。
首先把牛肉切成半个巴掌大小的块状，然后往肉上撒盐和黑椒粉，两面抹匀，腌制一两小时。剁碎蒜头和洋葱放入
热锅，待油融化后摆好牛排，慢火煎至你想要的成熟度。制作完成后将牛排装盘，用锅中的余油把蒜头碎和洋葱碎爆香。
随后把洋葱碎炒软后，加入一碗清水煮开，放几滴酒，再用盐和黑椒粉调味，慢慢炒匀成稠汁状。将熬好的黑椒汁淋到牛排上即可`
		resp, err := gptCaller.Call(ctx, fmt.Sprintf("{%v}", text))
		if err != nil {
			logrus.Errorf("gpt caller call fail, err = %v", err)
			return
		}
		logrus.Info(resp)
	}
	{
		text := `今天天气很不错，也不会下雨，也没有多云，阳光明媚`
		resp, err := gptCaller.Call(ctx, fmt.Sprintf("{%v}", text))
		if err != nil {
			logrus.Errorf("gpt caller call fail, err = %v", err)
			return
		}
		logrus.Info(resp)
	}
}

func promptJSON(ctx context.Context) {
	systemText := `你是一个文本总结器，可以把{}括起来的内容总结并提炼出对应的关键字。
你需要以JSON形式返回提炼文本、文本中涉及到的关键字列表，对应的Key分别为summarizeText、keywords`
	gptCaller, err := llm_caller.NewGptLLMCaller(ctx, systemText, 0.0, 512)
	if err != nil {
		logrus.Errorf("gpt caller init fail, err = %v,", err)
		return
	}
	{
		text := `大模型问世仅仅一年，就迎来了免费时代。
在字节跳动的大模型喊出比行业便宜99.3%之后，其“行业最低价”正在被不断打破。近日，阿里云、百度等巨头公司纷纷宣布将大模型价格降到更低，百度甚至直接宣布免费。科大讯飞更是宣布，其讯飞星火API能力正式免费开放。其中，讯飞星火Lite API永久免费开放。
多位业内人士在接受《中国经营报》记者采访时表示，尽管宣布降价的模型产品众多，但真正大规模、高性能并且支持高并发的大模型推理仍然需要收费，降价幅度有限。不过，可以明确的是，这场关于大模型免费使用的价格战才刚刚开始。也许很快人们就会看到，在大模型这场比赛中，大多数的大模型公司都将被淘汰。
没有最低，只有更低
5月22日，腾讯云正式对外宣布了全新的大模型升级方案，其核心模型之一的混元-lite模型价格策略发生重要调整，由原先的0.008元/千tokens调整为全面免费使用。同时，最高配置的万亿参数模型混元-pro也进行了价格优化，从原有的0.1元/千tokens下降至0.03元/千tokens。
短短数日，百度、阿里云、腾讯、字节跳动、科大讯飞等各家大模型厂商均已投入到这场“价格战”中。而国内这场“价格战”最早可追溯到5月6日，幻方量化旗下DeepSeek发布第二代MoE（专家模型）DeepSeek-V2，API定价为每百万tokens输入1元、输出2元（32K上下文），价格为GPT-4-Turbo的近百分之一。`
		resp, err := gptCaller.Call(ctx, fmt.Sprintf("{%v}", text))
		if err != nil {
			logrus.Errorf("gpt caller call fail, err = %v", err)
			return
		}
		logrus.Info(resp)
	}

	{
		text := `百度文心大模型5.0将在2025年发布
据知情人士消息，百度或将于2025年百度世界大会期间发布新一代文心大模型5.0。目前，文心大模型最新版本为4.0版本，具备理解、生成、逻辑和记忆四大核心能力。公开信息显示，百度文心大模型于2019年3月首发，同年7月迭代至2.0版本，2021年7月发布3.0版本，2023年10月升级至4.0版本。截至目前，百度方面暂未就此传闻给出回应。
OpenAI宣布已启动下一代前沿模型训练
当地时间周二，OpenAI发布公告称董事会成立了一个负责把控AI开发方向的安全委员会。同时，OpenAI还在公告中表示，近些日子已经开始训练公司的“下一代前沿模型”，该新模型预计带来更高水平的能力，将成为聊天机器人、智能助手、搜索引擎和图像生成器在内的各类人工智能产品的核心引擎，助力Open AI实现“通用人工智能”(AGI)的目标。`
		resp, err := gptCaller.Call(ctx, fmt.Sprintf("{%v}", text))
		if err != nil {
			logrus.Errorf("gpt caller call fail, err = %v", err)
			return
		}
		logrus.Info(resp)
	}
}

func promptDelimiter(ctx context.Context) {
	systemText := `你是一个文本总结器，可以把{}扩起来的内容总结为20个字内的概要`
	gptCaller, err := llm_caller.NewGptLLMCaller(ctx, systemText, 0.0, 512)
	if err != nil {
		logrus.Errorf("gpt caller init fail, err = %v,", err)
		return
	}

	{
		text := `大模型问世仅仅一年，就迎来了免费时代。
在字节跳动的大模型喊出比行业便宜99.3%之后，其“行业最低价”正在被不断打破。近日，阿里云、百度等巨头公司纷纷宣布将大模型价格降到更低，百度甚至直接宣布免费。科大讯飞更是宣布，其讯飞星火API能力正式免费开放。其中，讯飞星火Lite API永久免费开放。
多位业内人士在接受《中国经营报》记者采访时表示，尽管宣布降价的模型产品众多，但真正大规模、高性能并且支持高并发的大模型推理仍然需要收费，降价幅度有限。不过，可以明确的是，这场关于大模型免费使用的价格战才刚刚开始。也许很快人们就会看到，在大模型这场比赛中，大多数的大模型公司都将被淘汰。
没有最低，只有更低
5月22日，腾讯云正式对外宣布了全新的大模型升级方案，其核心模型之一的混元-lite模型价格策略发生重要调整，由原先的0.008元/千tokens调整为全面免费使用。同时，最高配置的万亿参数模型混元-pro也进行了价格优化，从原有的0.1元/千tokens下降至0.03元/千tokens。
短短数日，百度、阿里云、腾讯、字节跳动、科大讯飞等各家大模型厂商均已投入到这场“价格战”中。而国内这场“价格战”最早可追溯到5月6日，幻方量化旗下DeepSeek发布第二代MoE（专家模型）DeepSeek-V2，API定价为每百万tokens输入1元、输出2元（32K上下文），价格为GPT-4-Turbo的近百分之一。`
		resp, err := gptCaller.Call(ctx, fmt.Sprintf("{%v}", text))
		if err != nil {
			logrus.Errorf("gpt caller call fail, err = %v", err)
			return
		}
		logrus.Infof("-原始文本:%v\n-总结文本:%v", text, resp)
	}

	{
		text := `百度文心大模型5.0将在2025年发布
据知情人士消息，百度或将于2025年百度世界大会期间发布新一代文心大模型5.0。目前，文心大模型最新版本为4.0版本，具备理解、生成、逻辑和记忆四大核心能力。公开信息显示，百度文心大模型于2019年3月首发，同年7月迭代至2.0版本，2021年7月发布3.0版本，2023年10月升级至4.0版本。截至目前，百度方面暂未就此传闻给出回应。
OpenAI宣布已启动下一代前沿模型训练
当地时间周二，OpenAI发布公告称董事会成立了一个负责把控AI开发方向的安全委员会。同时，OpenAI还在公告中表示，近些日子已经开始训练公司的“下一代前沿模型”，该新模型预计带来更高水平的能力，将成为聊天机器人、智能助手、搜索引擎和图像生成器在内的各类人工智能产品的核心引擎，助力Open AI实现“通用人工智能”(AGI)的目标。`
		resp, err := gptCaller.Call(ctx, fmt.Sprintf("{%v}", text))
		if err != nil {
			logrus.Errorf("gpt caller call fail, err = %v", err)
			return
		}
		logrus.Infof("-原始文本:%v\n-总结文本:%v", text, resp)
	}
}
