package configs

import (
	"github.com/spf13/viper"
)

var cfg *conf

type conf struct {
	ApiKey               string `mapstructure:"API_KEY"`
	WebServerPort        string `mapstructure:"WEB_SERVER_PORT"`
	RedisAddress         string `mapstructure:"REDIS_ADDRESS"`
	RedisPassword        string `mapstructure:"REDIS_PASSWORD"`
	MaxRequestsPerIp     int    `mapstructure:"MAX_REQUESTS_PER_IP"`
	MaxRequestsPerApiKey int    `mapstructure:"MAX_REQUESTS_PER_API_KEY"`
	TimeInterval         int    `mapstructure:"TIME_INTERVAL_MINUTES"`
	BlockDuration        int    `mapstructure:"BLOCK_DURATION_MINUTES"`
}

func LoadConfig(path string) (*conf, error) {
	viper.SetConfigFile(path + "/.env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	return cfg, err
}
