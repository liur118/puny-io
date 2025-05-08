package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
)

type Config struct {
	Port      string            `yaml:"port" mapstructure:"port"`
	Storage   string            `yaml:"storage" mapstructure:"storage"`
	Users     map[string]string `yaml:"users" mapstructure:"users"`
	JwtSecret string            `yaml:"jwt_secret" mapstructure:"jwt_secret"`
	Host      string            `yaml:"host" mapstructure:"host"`
}

func (c *Config) Validate() error {
	if c.Port == "" {
		return fmt.Errorf("port is required")
	}
	if c.Storage == "" {
		return fmt.Errorf("storage path is required")
	}
	if c.JwtSecret == "" {
		return fmt.Errorf("jwt secret is required")
	}
	if len(c.Users) == 0 {
		return fmt.Errorf("at least one user is required")
	}
	return nil
}

func LoadConfig() (*Config, error) {
	pwd, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("failed to get executable path: %w", err)
	}
	viper.SetDefault("DecoderConfig", &mapstructure.DecoderConfig{
		TagName: "yaml",
	})

	configDir := filepath.Dir(pwd)

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(fmt.Sprintf("%s/conf", configDir))
	viper.AddConfigPath(fmt.Sprintf("%s/../conf", configDir))

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("使用的配置文件: %s\n", viper.ConfigFileUsed())
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	fmt.Printf("使用的配置文件: %s\n", viper.ConfigFileUsed())

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// 设置默认值
	if config.Port == "" {
		config.Port = "8080"
	}
	if config.Host == "" {
		config.Host = "http://localhost:8080"
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &config, nil
}
