package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DaseinBot DaseinBotConfig `mapstructure:"dasein_bot"`
	OpenAI    OpenAIConfig    `mapstructure:"openai"`
	MongoDB   MongoDBConfig   `mapstructure:"mongodb"`
	Redis     RedisConfig     `mapstructure:"redis"`
}

type DaseinBotConfig struct {
	Token              string   `mapstructure:"token"`
	CheckIfUserAllowed bool     `mapstructure:"check_if_user_allowed"`
	AllowedUsers       []int64  `mapstructure:"allowed_users"`
	Debug              bool     `mapstructure:"debug"`
	Questions          []string `mapstructure:"questions"`
	AnswerMaxLength    int      `mapstructure:"answer_max_length"`
	Start              string   `mapstructure:"start"`
	Help               string   `mapstructure:"help"`
	About              string   `mapstructure:"about"`
	OpenAIUserLimit    int      `mapstructure:"openai_user_limit"`
	DaysLeftMessage    string   `mapstructure:"days_left_message"`
}

type OpenAIConfig struct {
	APIKey           string `mapstructure:"api_key"`
	Model            string `mapstructure:"model"`
	Dummy            bool   `mapstructure:"dummy"`
	DeveloperMessage string `mapstructure:"developer_message"`
}

type MongoDBConfig struct {
	URI            string        `mapstructure:"uri"`
	Username       string        `mapstructure:"username"`
	Password       string        `mapstructure:"password"`
	Database       string        `mapstructure:"database"`
	Collection     string        `mapstructure:"collection"`
	ConnectTimeout time.Duration `mapstructure:"connect_timeout"`
}

type RedisConfig struct {
	Host     string        `mapstructure:"host"`
	Port     int           `mapstructure:"port"`
	Username string        `mapstructure:"username"`
	Password string        `mapstructure:"password"`
	DB       int           `mapstructure:"db"`
	TTL      time.Duration `mapstructure:"ttl"`
}

func NewConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)

	if err := viper.BindEnv("dasein_bot.token", "TELEGRAM_BOT_TOKEN"); err != nil {
		return nil, fmt.Errorf("error binding env var: %w", err)
	}
	if err := viper.BindEnv("openai.api_key", "OPENAI_API_KEY"); err != nil {
		return nil, fmt.Errorf("error binding env var: %w", err)
	}
	if err := viper.BindEnv("mongodb.password", "MONGO_PASSWORD"); err != nil {
		return nil, fmt.Errorf("error binding env var: %w", err)
	}
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %w", err)
	}
	return &cfg, nil
}
