package config

import "os"

type JWTConfig struct {
	Secret string
}

type Config struct {
	JWT JWTConfig
}

var JWT = JWTConfig{

	Secret: "your-secret-key", // Replace with your actual secret key

}

func NewConfig() *Config {
	return &Config{
		JWT: JWTConfig{
			Secret: getEnvOrDefault("JWT_SECRET", "jjdklfjajs;dlhgdsha;lkjsdlfasjkl"),
		},
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
