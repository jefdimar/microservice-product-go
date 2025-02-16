package config

import "fmt"

func (c *Config) Validate() error {
	if c.MongoDBName == "" {
		return fmt.Errorf("MONGO_DB_NAME is required")
	}

	if c.ServerAddress == "" {
		return fmt.Errorf("SERVER_ADDRESS is required")
	}

	return nil
}
