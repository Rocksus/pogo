package configs

import "os"

// New function returns a new config function filled with environment variables.
func New() *Config {
	return &Config{
		Chat: ChatConfig{
			MasterID:           getEnv("MASTER_ID", ""),
			ChannelAccessToken: getEnv("ACCESS_TOKEN", ""),
			ChannelSecret:      getEnv("CHANNEL_SECRET", ""),
		},
		Interpretor: InterpretorConfig{
			Token: getEnv("WIT_TOKEN_KEY", ""),
		},
		Port: getEnv("PORT", "8080"),
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
