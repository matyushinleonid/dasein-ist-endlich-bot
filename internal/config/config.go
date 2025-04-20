package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Token              string   `mapstructure:"token"`
	CheckIfUserAllowed bool     `mapstructure:"check_if_user_allowed"`
	AllowedUsers       []int64  `mapstructure:"allowed_users"`
	Questions          []string `mapstructure:"questions"`

	OpenAI OpenAIConfig `mapstructure:"openai"`
}

type OpenAIConfig struct {
	APIKey           string `mapstructure:"api_key"`
	Model            string `mapstructure:"model"`
	Dummy            bool   `mapstructure:"dummy"`
	DeveloperMessage string `mapstructure:"developer_message"`
}

func NewConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)

	if err := viper.BindEnv("openai.api_key", "OPENAI_API_KEY"); err != nil {
		return nil, fmt.Errorf("error binding env var: %w", err)
	}
	if err := viper.BindEnv("token", "TELEGRAM_BOT_TOKEN"); err != nil {
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
