package conf

import (
	"github.com/spf13/viper"
	"log"
)

var LLMHubConfig *Config

type OpenaiConfig struct {
	Key  string `json:"key"`
	Host string `json:"host"`
}

type RedisConfig struct {
	Url string `json:"url"`
}
type Config struct {
	Openai OpenaiConfig `json:"openai"`
	Redis  RedisConfig  `json:"redis"`
}

func Init() error {
	config := &Config{}
	vip := viper.New()
	vip.AddConfigPath("./conf")
	vip.AddConfigPath("../conf")
	vip.AddConfigPath("../../conf")
	vip.SetConfigName("config_private")
	vip.SetConfigType("yaml")
	err := vip.ReadInConfig()
	if err != nil {
		log.Printf("[config.Init] err = %v", err)
		return err
	}
	if err := vip.Unmarshal(config); err != nil {
		log.Printf("[config.Init] err = %v", err)
		return err
	}
	log.Printf("[config.Init] 初始化配置成功,config=%v", config)
	LLMHubConfig = config
	return nil
}
