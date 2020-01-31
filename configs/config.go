package configs

import "os"

// New function returns a new config function filled with environment variables.
func New() *Config {
	return &Config{
		Line: LineConfig{
			MasterID:           getEnv("MASTER_ID", ""),
			ChannelAccessToken: getEnv("ACCESS_TOKEN", ""),
			ChannelSecret:      getEnv("CHANNEL_SECRET", ""),
		},
		Wit: WitConfig{
			Token: getEnv("WIT_TOKEN_KEY", ""),
		},
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
