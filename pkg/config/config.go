package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	// Server
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`

	// MongoDB
	MongoHost     string `mapstructure:"MONGO_HOST"`
	MongoPort     string `mapstructure:"MONGO_PORT"`
	MongoUser     string `mapstructure:"MONGO_USER"`
	MongoPassword string `mapstructure:"MONGO_PASSWORD"`
	MongoDBName   string `mapstructure:"MONGO_DB_NAME"`

	// Redis
	RedisHost     string `mapstructure:"REDIS_HOST"`
	RedisPort     string `mapstructure:"REDIS_PORT"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AutomaticEnv()

	config := &Config{}

	// Set defaults
	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		// It's okay if config file doesn't exist
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	err := viper.Unmarshal(config)
	return config, err
}

func setDefaults() {
	viper.SetDefault("SERVER_ADDRESS", ":8001")
	viper.SetDefault("MONGO_HOST", "localhost")
	viper.SetDefault("MONGO_PORT", "27017")
	viper.SetDefault("MONGO_USER", "")
	viper.SetDefault("MONGO_PASSWORD", "")
	viper.SetDefault("MONGO_DB_NAME", "products")
	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_PORT", "6379")
	viper.SetDefault("REDIS_PASSWORD", "")
}
